package helper

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/rs/zerolog/log"
)

var ExcludedDirs = []string{
	"templates", ".github", "docs",
	".vscode", "tools", ".devcontainer",
	"testdata",
}

// WalkMode specifies the mode for walking charts
type WalkMode int

const (
	// SyncMode processes charts sequentially
	SyncMode WalkMode = iota
	// AsyncMode processes charts concurrently
	AsyncMode
)

type ActionFunc func(string, string) error

func getWalkDirFunc(action ActionFunc, bump string, mode WalkMode, wg *sync.WaitGroup) fs.WalkDirFunc {
	return func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && slices.Contains(ExcludedDirs, info.Name()) {
			return filepath.SkipDir
		}

		// Check if the current file is Chart.yaml
		if info.Name() == "Chart.yaml" {
			switch mode {
			case SyncMode:
				// Process charts sequentially
				if err := action(path, bump); err != nil {
					log.Fatal().Err(err).Msgf("Error executing action on Chart.yaml at [%s]", path)
				}
			case AsyncMode:
				// Process charts concurrently
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					if err := action(path, bump); err != nil {
						log.Fatal().Err(err).Msgf("Error executing action on Chart.yaml at [%s]", path)
					}
				}(path)
			default:
				return fmt.Errorf("invalid mode: %d", mode)
			}

			// Stop processing the current directory after finding Chart.yaml
			return filepath.SkipDir
		}

		return nil
	}
}

func UniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us

}
