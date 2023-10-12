package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/terminal"
	"github.com/Masterminds/sprig/v3"
)

// appConfig is the global application configuration.
var appConfig = config.GetAppConfig()

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

// injectCustomToFileFuncPaths takes a file path, its content, and a flag indicating whether the output
// is intended for the terminal. It returns a modified version of the file content with the two first hidden
// parameter of the 'toFile' functions injected.
//
// If 'outputInTerminal' is true, an error is returned because the 'toFile' function is not supported in
// terminal output mode; it should be used with 'cloney dry-run -o <output_directory>' instead.
func injectCustomToFileFuncPaths(templateDir, filePath, fileContent string, outputInTerminal bool) (string, error) {
	// Split the template content into lines for processing.
	fileLines := strings.Split(fileContent, "\n")

	// Iterate over each line in the template content.
	for index, line := range fileLines {
		// Inject the "hidden" parameters.
		regex := regexp.MustCompile(`{{-? ?toFile`)
		fileDir := filepath.Dir(filePath)

		// If on Windows, replace backslashes with forward slashes.
		if os.PathSeparator == '\\' {
			templateDir = filepath.ToSlash(templateDir)
			fileDir = filepath.ToSlash(fileDir)
		}
		newLine := regex.ReplaceAllString(line, fmt.Sprintf("{{- toFile \"%s\" \"%s\"", templateDir, fileDir))

		// If 'outputInTerminal' is true, return an error as "toFile" is not supported in terminal output mode.
		// if outputInTerminal {
		if newLine != line && outputInTerminal {
			return "", fmt.Errorf("the 'toFile' function is not supported when outputting the result to the terminal. Use 'cloney dry-run -o <output_directory>' instead")
		}

		// Replace the line in the file content.
		fileLines[index] = newLine
	}

	// Reconstruct the modified file content.
	newContent := strings.Join(fileLines, "\n")
	return newContent, nil
}

// FillDirectory processes template files in a source directory, replacing placeholders with variables.
func (t *TemplateFiller) FillDirectory(src string, ignorePaths []string, outputInTerminal bool) error {
	// Get a list of all files in the specified directory, considering ignore options.
	filePaths, err := GetAllFilePaths(src, ignorePaths)
	if err != nil {
		return fmt.Errorf("error obtaining file paths in directory %s: %w", src, err)
	}

	// Create a template and add custom functions.
	tmpl := template.New("")
	tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.Funcs(CustomTxtFuncMap(tmpl))

	// Create a map to hold the file contents.
	fileContents := make(map[string]string)

	// Iterate over each file in the directory and read the content.
	for _, filePath := range filePaths {
		// Read the content of the file.
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		// Get a new version of the file content with the first hidden parameter of the 'toFile' function injected.
		fileContent, err := injectCustomToFileFuncPaths(src, filePath, string(fileBytes), outputInTerminal)
		if err != nil {
			return err
		}

		fileContents[filePath] = fileContent
	}

	// Parse all file contents into the template.
	for filePath, fileContent := range fileContents {
		_, err = tmpl.New(filePath).Parse(fileContent)
		if err != nil {
			return fmt.Errorf("error parsing template for file %s: %w", filePath, err)
		}
	}

	// Execute the templates for each file.
	for _, filePath := range filePaths {
		var resultBuffer bytes.Buffer
		err = tmpl.ExecuteTemplate(&resultBuffer, filePath, t.Variables)
		if err != nil {
			return fmt.Errorf("error executing template for file %s: %w", filePath, err)
		}

		// If the 'outputInTerminal' parameter is set, output the result to the terminal.
		if outputInTerminal {
			if !strings.HasPrefix(filepath.Base(filePath), appConfig.IgnorePrefix) {
				terminal.Message(fmt.Sprintf("\n--- File (%s)\n%s\n", terminal.Blue(filePath), resultBuffer.String()))
			}
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
