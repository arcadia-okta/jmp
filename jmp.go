package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

const FAVOURITE_MULTIPLIER int = 10

type config struct {
	ProjectsRoot string   `json:"projectsRoot"`
	SearchDepth  int      `json:"searchDepth"`
	Favourites   []string `json:"favourites"`
}

func isGitRepo(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}
	for _, file := range files {
		if file.Name() == ".git" && file.IsDir() {
			return true, nil
		}
	}
	return false, nil
}

func findGitRepos(start string, maxDepth int) ([]string, error) {
	var dirs []string = []string{}

	if maxDepth == -1 {
		return dirs, nil
	}

	files, err := ioutil.ReadDir(start)
	if err != nil {
		return dirs, err
	}

	for _, file := range files {
		pathToFile := path.Join(start, file.Name())
		if file.IsDir() {
			isRepo, err := isGitRepo(pathToFile)
			if err != nil {
				return dirs, err
			} else if isRepo {
				dirs = append(dirs, pathToFile)
			} else {
				repos, err := findGitRepos(pathToFile, maxDepth-1)
				if err != nil {
					return dirs, err
				}
				dirs = append(dirs, repos...)
			}
		}
	}

	return dirs, nil
}

func isFavourite(repo string, faves []string) bool {
	for _, fave := range faves {
		if repo == fave {
			return true
		}
	}
	return false
}

func fuzzyFindGitRepo(cfg config, target string) (string, error) {
	gitRepos, err := findGitRepos(os.ExpandEnv(cfg.ProjectsRoot), cfg.SearchDepth)
	if err != nil {
		return "", nil
	}

	bestPath := ""
	bestScore := -1

	for _, repoPath := range gitRepos {
		_, directoryName := path.Split(repoPath)
		score := fuzzy.RankMatch(target, directoryName)
		if isFavourite(directoryName, cfg.Favourites) && score >= 1 {
			score *= FAVOURITE_MULTIPLIER
		}
		if score > bestScore {
			bestPath = repoPath
			bestScore = score
		}
	}

	return bestPath, nil
}

func usage() string {
	return `jmp <repository>
Jump to a git repository.`
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		os.Exit(1)
	}
	targetRepo := os.Args[1]

	cfg := config{}
	configPath := path.Join(os.ExpandEnv("$HOME"), ".jmp", "cfg.json")
	cfgFile, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open config file %s, reason: %s.\n", configPath, err.Error())
		os.Exit(1)
	}
	decoder := json.NewDecoder(cfgFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load configuration in %s, reason: %s.\n", configPath, err.Error())
		os.Exit(1)
	}

	path, err := fuzzyFindGitRepo(cfg, targetRepo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't find a repository matching %s. Error = %s\n", targetRepo, err.Error())
		os.Exit(1)
	} else if path == "" {
		fmt.Fprintf(os.Stderr, "Couldn't find a repository matching %s.\n", targetRepo)
		os.Exit(1)
	}

	fmt.Println(path)
}
