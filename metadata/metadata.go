package metadata

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// CloneyVariable is the struct that represents a variable in a Cloney template repository.
type CloneyVariable struct {
	// Name is the name of the variable.
	Name string `yaml:"name"`

	// Description is the description of the variable.
	Description string `yaml:"description"`

	// Default is the default value of the variable.
	Default string `yaml:"default"`

	// Type is the type of the variable.
	Type string `yaml:"type"`
}

// CloneyMetadata is the struct that represents the metadata file of a Cloney template repository.
type CloneyMetadata struct {
	// Name is the name of the template repository.
	Name string `yaml:"name" validate:"required"`

	// Description is the description of the template repository.
	Description string `yaml:"description"`

	// Version is the version of the template repository.
	Version string `yaml:"version" validate:"required" validate:"semver"`

	// Authors is the list of authors of the template repository.
	Authors []string `yaml:"authors"`

	// License is the license of the template repository.
	License string `yaml:"license"`

	// Variables is the list of variables of the template repository.
	Variables []CloneyVariable `yaml:"variables"`
}

// NewCloneyMetadata creates a new CloneyMetadata struct from a YAML string.
func NewCloneyMetadata(stringYAMLContent string) (*CloneyMetadata, error) {
	// Parse YAML.
	var metadata CloneyMetadata
	err := yaml.Unmarshal([]byte(stringYAMLContent), &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// PrettyPrint prints the Cloney template repository metadata in a pretty way.
func (m *CloneyMetadata) PrettyPrint() {
	// Print basic information.
	fmt.Printf("%s\n\n", m.Name)
	if m.Description != "" {
		fmt.Printf("Description: %s\n\n", m.Description)
	} else {
		fmt.Print("No description provided.\n\n")
	}
	fmt.Printf("Version: %s\n", m.Version)
	if len(m.Authors) > 0 {
		fmt.Printf("Authors: %s\n", m.Authors)
	}
	if m.License != "" {
		fmt.Printf("License: %s\n", m.License)
	}

	// Print variables.
	if len(m.Variables) == 0 {
		fmt.Print("\nNo variables provided.\n")
		return
	}
	fmt.Print("\nVariables:\n\n")
	variablesTable := tablewriter.NewWriter(os.Stdout)
	variablesTable.SetHeader([]string{"Name", "Description", "Type", "Default"})
	for _, variable := range m.Variables {
		variablesTable.Append(
			[]string{variable.Name, variable.Description, variable.Type, variable.Default},
		)
	}
	variablesTable.Render()
}
