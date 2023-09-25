package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetAllFilePaths returns a list of all file paths within a directory and its subdirectories,
// with options to specify directories and files to ignore.
func GetAllFilePaths(directoryPath string, ignorePaths []string) ([]string, error) {
	var filePaths []string

	// Walk the directory and its subdirectories.
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}

		// Check if the path should be ignored.
		var ignore bool
		for _, ignorePath := range ignorePaths {
			absPath, _ := filepath.Abs(path)
			absDirPath, _ := filepath.Abs(directoryPath)
			absIgnorePath := filepath.Join(absDirPath, ignorePath)
			if strings.HasPrefix(absPath, absIgnorePath) {
				ignore = true
				break
			}
		}

		if info.IsDir() {
			// Check if the directory should be ignored.
			if ignore {
				return filepath.SkipDir
			}
		} else {
			// Check if the file should be ignored.
			if ignore {
				return nil
			}
			// Add the file path to the list.
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", directoryPath, err)
	}

	return filePaths, nil
}
