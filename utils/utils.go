package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
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

// IgnorePathOptions is a struct used to configure the GetAllFilePaths function.
type IgnorePathOptions struct {
	// IgnoreFiles is a list of file names to ignore.
	IgnoreFiles []string

	// IgnoreDirectories is a list of directory names to ignore.
	IgnoreDirectories []string
}

// GetAllFilePaths returns a list of all file paths within a directory and its subdirectories,
// while allowing you to specify directories to ignore.
func GetAllFilePaths(directoryPath string, ignoreOptions IgnorePathOptions) ([]string, error) {
	var filePaths []string

	// Walk the directory and its subdirectories.
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if info.IsDir() {
			// Check if the directory should be ignored.
			if ListContainsString(ignoreOptions.IgnoreDirectories, info.Name()) {
				return filepath.SkipDir
			}
		} else {
			// Check if the file should be ignored.
			if ListContainsString(ignoreOptions.IgnoreFiles, info.Name()) {
				return nil
			}
			// Add file path to the list.
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", directoryPath, err)
	}

	return filePaths, nil
}

// Colors in the terminal.
var Green = color.New(color.FgGreen).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()
var Blue = color.New(color.FgBlue).SprintFunc()

// Messages.

// OKMessage prints a message with a green [OK] prefix.
func OKMessage(message string) {
	fmt.Printf("%s %s.\n", Green("[OK]"), message)
}

// WarningMessage prints a message with a yellow [Warning] prefix.
func WarningMessage(message string) {
	fmt.Printf("%s %s.\n", Yellow("[Warning]"), message)
}

// ErrorMessage prints a message with a red [Error] prefix.
func ErrorMessage(message string, err error) {
	fmt.Printf("%s %s: %v.\n", Red("[Error]"), message, err)
}

// CommandExampleWithExplanation return a string with the command example and explanation
// using an yellow color to separate the example from the explanation.
func CommandExampleWithExplanation(example, explanation string) string {
	return fmt.Sprintf("%s %s %s", example, Yellow("==>"), explanation)
}
