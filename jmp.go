package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func findDir(start, target string, maxDepth int) (string, error) {
	if maxDepth == -1 {
		return "", nil
	}

	files, err := ioutil.ReadDir(start)
	if err != nil {
		return "", err
	}

	var dirs []string = make([]string, 1)
	for _, file := range files {
		if file.IsDir() {
			// TODO - Use a fuzzy match.
			if file.Name() == target {
				return path.Join(start, file.Name()), nil
			} else {
				dirs = append(dirs, file.Name())
			}
		}
	}

	for _, dir := range dirs {
		nextDir := path.Join(start, dir)
		found, err := findDir(nextDir, target, maxDepth-1)
		if err == nil && found != "" {
			return found, nil
		}
	}

	return "", nil
}

func main() {
	start := "/Users/arcadia.rose/Code"
	targetRepo := os.Args[1]

	path, err := findDir(start, targetRepo, 5)
	found := err == nil && path != ""
	if found {
		fmt.Printf("Found %s\n", path)
	} else {
		fmt.Printf("Couldn't find it. Error = %s\n", err.Error())
	}
}
