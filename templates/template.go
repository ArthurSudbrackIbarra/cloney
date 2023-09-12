package templates

import (
	"bytes"
	"fmt"
	"os"
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

// FillDirectory recursively processes all files in a directory with the variables from the TemplateFiller.
func (t *TemplateFiller) FillDirectory(directoryPath string) error {
	// Get a list of all files in the specified directory, excluding the ".git" directory.
	filePaths, err := utils.GetAllFilePaths(directoryPath, []string{
		".git",
	})
	if err != nil {
		return fmt.Errorf("error while obtaining file paths in directory %s: %w", directoryPath, err)
	}

	// Iterate over each file in the directory, applying the template to each file and saving the result in the same file.
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

		// Write the resulting content back to the file, overwriting the original file.
		err = os.WriteFile(filePath, resultBuffer.Bytes(), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error while writing file %s: %w", filePath, err)
		}
	}

	return nil
}
