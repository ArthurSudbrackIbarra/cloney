package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CopyDirectory copies a directory recursively
// with options to specify directories and files to ignore.
func CopyDirectory(src string, dest string, ignorePaths []string) error {
	// Get a list of all files in the specified directory, considering ignore options.
	filePaths, err := GetAllFilePaths(src, ignorePaths)
	if err != nil {
		return fmt.Errorf("error obtaining file paths in directory %s: %w", src, err)
	}

	// Iterate over each file in the directory and copy it.
	for _, filePath := range filePaths {
		// Get the relative path of the file.
		relativePath, err := filepath.Rel(src, filePath)
		if err != nil {
			return fmt.Errorf("error getting relative path of file %s: %w", filePath, err)
		}

		// Construct the destination path.
		destPath := filepath.Join(dest, relativePath)

		// Create the destination directory if it does not exist.
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", filepath.Dir(destPath), err)
		}

		// Read the file content.
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		// Write the file content to the destination.
		err = os.WriteFile(destPath, fileContent, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error copying file %s: %w", filePath, err)
		}
	}

	return nil
}

// ShouldIgnorePath determines whether a given file or directory path should be ignored
// based on a list of patterns within a specified base directory. It returns true if the
// path should be ignored according to any of the provided patterns, and false otherwise.
func ShouldIgnorePath(baseDirectory string, path string, ignorePaths []string) (bool, error) {
	for _, ignorePath := range ignorePaths {
		// Construct the full path for comparison.
		fullIgnorePath := filepath.Join(baseDirectory, ignorePath)
		match, err := filepath.Match(fullIgnorePath, path)
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

		// Skip the root directory.
		if path == directoryPath {
			return nil
		}

		// Check if the path should be ignored.
		ignore, err := ShouldIgnorePath(directoryPath, path, ignorePaths)
		if err != nil {
			return fmt.Errorf("error checking if path %s should be ignored: %w", path, err)
		}

		if info.IsDir() {
			// Check if the directory should be ignored and skip it if necessary.
			if ignore {
				return filepath.SkipDir
			}
		} else {
			// Check if the file should be ignored and skip it if necessary.
			if !ignore {
				filePaths = append(filePaths, path) // Add the file path to the list.
			}
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
// Additionally, it deletes files and directories starting with "_", which are files that should
// be processed by Cloney but not copied when the template is cloned.
func DeleteIgnoredFiles(directoryPath string, ignorePaths []string) error {
	// Walk the directory and its subdirectories.
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}

		// Skip the root directory.
		if path == directoryPath {
			return nil
		}

		// Check if the path starts with "_".
		var delete bool
		if strings.HasPrefix(filepath.Base(path), "_") {
			delete = true
		}

		// Check if the path should be ignored.
		if !delete {
			delete, err = ShouldIgnorePath(directoryPath, path, ignorePaths)
			if err != nil {
				return fmt.Errorf("error checking if path %s should be ignored: %w", path, err)
			}
		}

		// If the path should be ignored, delete it.
		if info.IsDir() && delete {
			err := os.RemoveAll(path)
			// Check if the error is due to the directory not existing.
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("error removing directory %s: %w", path, err)
			}
		} else if delete {
			err := os.Remove(path)
			// Check if the error is due to the file not existing.
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
