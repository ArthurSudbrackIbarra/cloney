package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

// ShouldIgnorePath determines whether a given file or directory path should be ignored
// based on a list of patterns within a specified base directory. It returns true if the
// path should be ignored according to any of the provided patterns, and false otherwise.
func ShouldIgnorePath(baseDirectory string, path string, ignorePaths []string) (bool, error) {
	for _, ignorePath := range ignorePaths {
		newIgnorePath := filepath.Join(baseDirectory, ignorePath)
		match, err := filepath.Match(newIgnorePath, path)
		if err != nil {
			return false, fmt.Errorf("error matching path %s: %w", path, err)
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

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
		ignore, err := ShouldIgnorePath(directoryPath, path, ignorePaths)
		if err != nil {
			return fmt.Errorf("error checking if path %s should be ignored: %w", path, err)
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

// DeleteIgnoredFiles recursively walks through a directory and its subdirectories
// specified by 'directoryPath'. It deletes files and directories that match any
// of the patterns in 'ignorePaths'.
func DeleteIgnoredFiles(directoryPath string, ignorePaths []string) error {
	// Walk the directory and its subdirectories
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}

		// Check if the path should be ignored.
		ignore, err := ShouldIgnorePath(directoryPath, path, ignorePaths)
		if err != nil {
			return fmt.Errorf("error checking if path %s should be ignored: %w", path, err)
		}

		// If the path should be ignored, delete it.
		if info.IsDir() && ignore {
			err := os.RemoveAll(path)
			// Check if error is because the directory does not exist.
			// If not, return the error.
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("error removing directory %s: %w", path, err)
			}
		} else if ignore {
			err := os.Remove(path)
			// Check if error is because the file does not exist.
			// If not, return the error.
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("error removing file %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory %s: %w", directoryPath, err)
	}

	return nil
}
