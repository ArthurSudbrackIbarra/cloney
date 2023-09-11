package template

import (
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// TemplateFiller is a struct that contains the variables to be used in the template.
type TemplateFiller struct {
	// Variables to be used in the template.
	Variables map[string]interface{}
}

// NewTemplateFiller creates a new TemplateFiller instance.
func NewTemplateFiller() *TemplateFiller {
	return &TemplateFiller{
		Variables: make(map[string]interface{}),
	}
}

// AddVariable adds a new variable to the TemplateFiller.
func (t *TemplateFiller) AddVariable(name string, value interface{}) {
	t.Variables[name] = value
}

// FillDirectory fills all the files in a directory with the variables in the TemplateFiller.
func (t *TemplateFiller) FillDirectory(directoryPath string) error {
	// Read directory files.
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return fmt.Errorf("Error reading directory: %w", err)
	}

	// Create the text template.
	template := template.New("cloneyTemplate").Funcs(sprig.TxtFuncMap())
}
