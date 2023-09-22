package templates

import (
	"fmt"
	"os"
	"path/filepath"

	basicoperations "github.com/ArthurSudbrackIbarra/cloney/basic-operations"
)

// IgnorePathOptions is a struct used to configure the GetAllFilePaths function.
type IgnorePathOptions struct {
	// IgnoreFiles is a list of file names to be ignored.
	IgnoreFiles []string

	// IgnoreDirectories is a list of directory names to be ignored.
	IgnoreDirectories []string
}

// GetAllFilePaths returns a list of all file paths within a directory and its subdirectories,
// with options to specify directories and files to ignore.
func GetAllFilePaths(directoryPath string, ignoreOptions IgnorePathOptions) ([]string, error) {
	var filePaths []string

	// Walk the directory and its subdirectories.
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if info.IsDir() {
			// Check if the directory should be ignored.
			if basicoperations.ListContainsString(ignoreOptions.IgnoreDirectories, info.Name()) {
				return filepath.SkipDir
			}
		} else {
			// Check if the file should be ignored.
			if basicoperations.ListContainsString(ignoreOptions.IgnoreFiles, info.Name()) {
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
