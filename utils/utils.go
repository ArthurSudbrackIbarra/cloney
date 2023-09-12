package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ListContainsString checks if a given list contains a specific string value.
func ListContainsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// GetAllFilePaths returns a list of all file paths within a directory and its subdirectories,
// while allowing you to specify directories to ignore.
func GetAllFilePaths(directoryPath string, ignoreDirectories []string) ([]string, error) {
	var filePaths []string

	// Walk the directory and its subdirectories.
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if !info.IsDir() {
			// Add file path to the list.
			filePaths = append(filePaths, path)
		} else {
			// Check if the directory should be ignored.
			if ListContainsString(ignoreDirectories, info.Name()) {
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePaths, nil
}
