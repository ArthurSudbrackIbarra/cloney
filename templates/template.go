package templates

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

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

// AddVariable adds a new variable to the TemplateFiller instance.
func (t *TemplateFiller) AddVariable(name string, value interface{}) {
	t.Variables[name] = value
}

// TemplateFillOptions is a struct used to configure the FillDirectory function.
type TemplateFillOptions struct {
	// SourceDirectoryPath is the path to the directory containing the template files.
	SourceDirectoryPath string

	// TargetDirectoryPath is the path to the directory where the filled template files will be saved.
	// If nil, the files will be saved in the same directory as the template files.
	TargetDirectoryPath *string

	// PrintMode, if true, causes the filled template files to be printed to a stdout instead of being saved to files.
	PrintMode bool

	// Stdout is the writer where the filled template files will be printed if PrintMode is enabled.
	Stdout io.Writer
}

// FillDirectory processes all files in a directory with variables from the TemplateFiller.
// If PrintMode is enabled, the filled template content is printed to Stdout instead of being saved to files.
func (t *TemplateFiller) FillDirectory(templateOptions TemplateFillOptions, ignoreOptions IgnorePathOptions) error {
	// If PrintMode is enabled, but Stdout is not defined, use os.Stdout.
	if templateOptions.PrintMode && templateOptions.Stdout == nil {
		templateOptions.Stdout = os.Stdout
	}

	// Get a list of all files in the specified directory.
	filePaths, err := GetAllFilePaths(templateOptions.SourceDirectoryPath, ignoreOptions)
	if err != nil {
		return fmt.Errorf("error obtaining file paths in directory %s: %w", templateOptions.SourceDirectoryPath, err)
	}

	// Iterate over each file in the directory, applying the template to each file.
	for _, filePath := range filePaths {
		// Read the content of the file.
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		// Create the template and add custom functions.
		tmpl := template.New("cloneyTemplate")
		tmpl.Funcs(sprig.TxtFuncMap()).Funcs(CustomTxtFuncMap(tmpl))

		// Parse the template.
		tmpl, err = tmpl.Parse(string(fileContent))
		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		// Execute the template, replacing placeholders with variables.
		var resultBuffer bytes.Buffer
		err = tmpl.Execute(&resultBuffer, t.Variables)
		if err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}

		// If PrintMode is enabled, print the result to Stdout and continue to the next file.
		if templateOptions.PrintMode {
			templateOptions.Stdout.Write(
				[]byte(fmt.Sprintf("\n----- File: %s\n%s\n", filePath, resultBuffer.String())),
			)
			continue
		}

		// Write the resulting content back to the file, overwriting the original file if TargetDirectoryPath is not defined.
		if templateOptions.TargetDirectoryPath == nil {
			err = os.WriteFile(filePath, resultBuffer.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing file %s: %w", filePath, err)
			}
		} else {
			// Calculate the path of the file relative to the source directory to preserve the directory structure.
			relativeFilePath, err := filepath.Rel(templateOptions.SourceDirectoryPath, filePath)
			if err != nil {
				return fmt.Errorf("error calculating relative file path: %w", err)
			}

			// Calculate the path of the file in the target directory.
			targetFilePath := filepath.Join(*templateOptions.TargetDirectoryPath, relativeFilePath)

			// If necessary, create the directory where the file will be saved.
			directory := filepath.Dir(targetFilePath)
			err = os.MkdirAll(directory, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %w", directory, err)
			}

			// Write the resulting content to the file.
			err = os.WriteFile(targetFilePath, resultBuffer.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing file %s: %w", targetFilePath, err)
			}
		}
	}

	return nil
}
