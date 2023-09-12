package template

import (
	"fmt"
	"os"
	"text/template"

	"github.com/ArthurSudbrackIbarra/cloney/utils"
	"github.com/Masterminds/sprig/v3"
)

// TemplateFiller is a struct that contains the variables to be used in the template.
type TemplateFiller struct {
	// Variables to be used in the template.
	Variables map[string]interface{}
}

// NewTemplateFiller creates a new TemplateFiller instance.
func NewTemplateFiller(variablesMap map[string]interface{}) *TemplateFiller {
	return &TemplateFiller{
		Variables: variablesMap,
	}
}

// AddVariable adds a new variable to the TemplateFiller.
func (t *TemplateFiller) AddVariable(name string, value interface{}) {
	t.Variables[name] = value
}

// FillDirectory recursively fills all the files in a directory with the variables in the TemplateFiller.
func (t *TemplateFiller) FillDirectory(directoryPath string) error {
	// Get all the files in the directory.
	// Ignore the .git directory.
	filePaths, err := utils.GetAllFilePaths(directoryPath, []string{
		".git",
	})
	if err != nil {
		return fmt.Errorf("error getting files in directory %s: %w", directoryPath, err)
	}

	fmt.Println(filePaths)

	// Create the text template.
	tmpl, err := template.New("cloneyTemplate").Funcs(sprig.TxtFuncMap()).ParseFiles(filePaths...)
	if err != nil {
		return fmt.Errorf("error creating template: %w", err)
	}

	// Fill the template.
	err = tmpl.Execute(os.Stdout, t.Variables)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil

}
