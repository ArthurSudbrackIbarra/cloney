package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ArthurSudbrackIbarra/cloney/utils"
	"github.com/Masterminds/sprig/v3"
)

// TemplateFiller is a struct used to fill Go templates with variables.
type TemplateFiller struct {
	// Variables to be injected into the template.
	Variables map[string]interface{}
}

// NewTemplateFiller creates a new TemplateFiller instance with the provided variables.
func NewTemplateFiller(variablesMap map[string]interface{}) *TemplateFiller {
	return &TemplateFiller{
		Variables: variablesMap,
	}
}

// AddVariable adds a new variable to the TemplateFiller.
func (t *TemplateFiller) AddVariable(name string, value interface{}) {
	t.Variables[name] = value
}

// TemplateFillOptions is a struct used to configure the FillDirectory function.
type TemplateFillOptions struct {
	// Path to the directory containing the template files.
	SourceDirectoryPath string

	// Path to the directory where the filled template files will be saved.
	// It is a pointer because it can be nil, in which case the files will be saved in the same directory as the template files.
	TargetDirectoryPath *string

	// If true, the filled template files will be printed to the terminal instead of being saved to the files.
	TerminalMode bool
}

// FillDirectory recursively processes all files in a directory with the variables from the TemplateFiller.
// If dryrunMode is enabled, the resulting content will be printed to the terminal instead of writing to the files.
func (t *TemplateFiller) FillDirectory(options TemplateFillOptions) error {
	// Get a list of all files in the specified directory, excluding the ".git" directory.
	filePaths, err := utils.GetAllFilePaths(options.SourceDirectoryPath, []string{
		".git",
	})
	if err != nil {
		return fmt.Errorf("error while obtaining file paths in directory %s: %w", options.SourceDirectoryPath, err)
	}

	// Iterate over each file in the directory, applying the template to each file.
	for _, filePath := range filePaths {
		// Read the content of the file.
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error while reading file %s: %w", filePath, err)
		}

		// Create the template and add the custom functions.
		tmpl := template.New("cloneyTemplate")
		tmpl.Funcs(sprig.TxtFuncMap()).Funcs(CustomTxtFuncMap(tmpl))

		// Parse the template.
		tmpl, err = tmpl.Parse(string(fileContent))
		if err != nil {
			return fmt.Errorf("error while parsing template: %w", err)
		}

		// Execute the template, replacing placeholders with variables.
		var resultBuffer bytes.Buffer
		err = tmpl.Execute(&resultBuffer, t.Variables)
		if err != nil {
			return fmt.Errorf("error while executing template: %w", err)
		}

		// If terminal mode is enabled, print the result to the terminal instead of writing it to the file.
		if options.TerminalMode {
			fmt.Printf("File: %s\n%s\n\n", filePath, resultBuffer.String())
			continue
		}

		// Write the resulting content back to the file, overwriting the original file, if TargetDirectoryPath is not defined.
		if options.TargetDirectoryPath == nil {
			err = os.WriteFile(filePath, resultBuffer.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error while writing file %s: %w", filePath, err)
			}
		} else {
			// Calculate the path of the file relative to the source directory.
			// This is done to preserve the directory structure when copying files.
			relativeFilePath, err := filepath.Rel(options.SourceDirectoryPath, filePath)
			if err != nil {
				return fmt.Errorf("error while calculating relative file path: %w", err)
			}

			// Calculate the path of the file in the target directory.
			targetFilePath := filepath.Join(*options.TargetDirectoryPath, relativeFilePath)

			// If necessary, create the directory where the file will be saved.
			directory := filepath.Dir(targetFilePath)
			err = os.MkdirAll(directory, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %w", directory, err)
			}

			// Write the resulting content to the file.
			err = os.WriteFile(targetFilePath, resultBuffer.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error while writing file %s: %w", targetFilePath, err)
			}
		}
	}

	return nil
}
