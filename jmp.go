package main

import (
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

func main() {
	start := "/Users/arcadia.rose/Code"
	targetRepo := os.Args[1]

	path, err := fuzzyFindGitRepo(start, targetRepo)
	found := err == nil && path != ""
	if found {
		fmt.Printf("Found %s\n", path)
	} else {
		fmt.Printf("Couldn't find it. Error = %s\n", err.Error())
	}
}
