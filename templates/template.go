package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ArthurSudbrackIbarra/cloney/terminal"
	"github.com/Masterminds/sprig/v3"
)

// TemplateFiller is a struct used for populating Go templates with variables.
type TemplateFiller struct {
	// Variables contains the variables to be injected into the template.
	Variables map[string]interface{}
}

// NewTemplateFiller creates a new TemplateFiller instance initialized with the provided variables.
func NewTemplateFiller(variablesMap map[string]interface{}) *TemplateFiller {
	return &TemplateFiller{
		Variables: variablesMap,
	}
}

// replaceCustomToFileFuncPaths replaces the paths of the "toFile" custom function in the template content.
// This function is necessary because the paths must be relative to the directory of the file being processed.
// The "toFile" function is defined in templates/custom_funcs.go.
// If 'outputInTerminal' is set to true, an error is returned as "toFile" is not supported in terminal output mode.
func replaceCustomToFileFuncPaths(filePath string, fileContent string, outputInTerminal bool) (string, error) {
	// Define a regular expression to match the "toFile" function and extract the path parameter.
	regex := regexp.MustCompile(`{{-?\s*toFile\s+"([^"]*)"`)

	// Split the template content into lines for processing.
	fileLines := strings.Split(fileContent, "\n")

	// Iterate over each line in the template content.
	for index, line := range fileLines {
		if strings.Contains(line, "toFile") {
			// If 'outputInTerminal' is true, return an error as "toFile" is not supported in terminal output mode.
			if outputInTerminal {
				return "", fmt.Errorf("the 'toFile' function is not supported when outputting the result to the terminal. Use 'cloney dry-run -o <output_directory>' instead")
			}
			// Extract the path parameter from the "toFile" function.
			matches := regex.FindStringSubmatch(line)
			if len(matches) != 2 {
				return "", fmt.Errorf("error parsing 'toFile' function path: %s", line)
			}
			pathParam := matches[1]

			// Check if the path is absolute, and if so, raise an error.
			if filepath.IsAbs(pathParam) {
				return "", fmt.Errorf("'toFile' function path must be relative: %s", pathParam)
			}

			// Calculate the new path of the file based on the template's location.
			newPath := filepath.Join(filepath.Dir(filePath), pathParam)
			newLine := strings.ReplaceAll(line, pathParam, newPath)

			// If on Windows, replace backslashes with forward slashes.
			if os.PathSeparator == '\\' {
				newLine = strings.ReplaceAll(newLine, "\\", "/")
			}

			// Replace the line in the file content.
			fileLines[index] = newLine
		}
	}

	// Reconstruct the modified file content.
	newContent := strings.Join(fileLines, "\n")
	return newContent, nil
}

// FillDirectory processes template files in a source directory, replacing placeholders with variables.
func (t *TemplateFiller) FillDirectory(src string, ignoreOptions IgnorePathOptions, outputInTerminal bool) error {
	// Get a list of all files in the specified directory, considering ignore options.
	filePaths, err := GetAllFilePaths(src, ignoreOptions)
	if err != nil {
		return fmt.Errorf("error obtaining file paths in directory %s: %w", src, err)
	}

	// Iterate over each file in the directory, applying the template to each file.
	for _, filePath := range filePaths {
		// Read the content of the file.
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		// Replace the paths of the "toFile" custom function in the template.
		fileContent, err := replaceCustomToFileFuncPaths(filePath, string(fileBytes), outputInTerminal)
		if err != nil {
			return err
		}

		// Create a template and add custom functions.
		tmpl := template.New("cloneyTemplate")
		tmpl.Funcs(sprig.TxtFuncMap()).Funcs(CustomTxtFuncMap(tmpl))

		// Parse the template.
		tmpl, err = tmpl.Parse(fileContent)
		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		// Execute the template, replacing placeholders with variables.
		var resultBuffer bytes.Buffer
		err = tmpl.Execute(&resultBuffer, t.Variables)
		if err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}

		// If the 'outputInTerminal' parameter is set, output the result to the terminal.
		if outputInTerminal {
			terminal.Message(fmt.Sprintf("\n----- File: %s\n%s\n", filePath, resultBuffer.String()))
		} else {
			// Write the result to the same file.
			err = os.WriteFile(filePath, resultBuffer.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// CreateFilledDirectory processes template files in a source directory and saves the filled files in a destination directory.
func (t *TemplateFiller) CreateFilledDirectory(src string, dest string, ignoreOptions IgnorePathOptions) error {
	// Get a list of all files in the specified directory, considering ignore options.
	filePaths, err := GetAllFilePaths(src, ignoreOptions)
	if err != nil {
		return fmt.Errorf("error obtaining file paths in directory %s: %w", src, err)
	}

	// Iterate over each file in the directory, applying the template to each file.
	for _, filePath := range filePaths {
		// Read the content of the file.
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		// Calculate the path of the file relative to the source directory to preserve the directory structure.
		relativeFilePath, err := filepath.Rel(src, filePath)
		if err != nil {
			return fmt.Errorf("error calculating relative file path: %w", err)
		}

		// Calculate the path of the file in the target directory.
		targetFilePath := filepath.Join(dest, relativeFilePath)

		// If necessary, create the directory where the file will be saved.
		directory := filepath.Dir(targetFilePath)
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", directory, err)
		}

		// Replace the paths of the "toFile" custom function in the template.
		fileContent, err := replaceCustomToFileFuncPaths(targetFilePath, string(fileBytes), false)
		if err != nil {
			return err
		}

		// Create a template and add custom functions.
		tmpl := template.New("cloneyTemplate")
		tmpl.Funcs(sprig.TxtFuncMap()).Funcs(CustomTxtFuncMap(tmpl))

		// Parse the template.
		tmpl, err = tmpl.Parse(fileContent)
		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		// Execute the template, replacing placeholders with variables.
		var resultBuffer bytes.Buffer
		err = tmpl.Execute(&resultBuffer, t.Variables)
		if err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}

		// Write the resulting content to the file in the target directory.
		err = os.WriteFile(targetFilePath, resultBuffer.Bytes(), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error writing file %s: %w", targetFilePath, err)
		}
	}

	return nil
}
