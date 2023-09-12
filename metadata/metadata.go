package metadata

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

// CloneyVariable represents a variable in a Cloney template repository.
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

// CloneyMetadata represents the metadata file of a Cloney template repository.
type CloneyMetadata struct {
	// Name is the name of the template repository.
	Name string `yaml:"name" validate:"required"`

	// Description is the description of the template repository.
	Description string `yaml:"description"`

	// Version is the version of the template repository.
	Version string `yaml:"version" validate:"required semver"`

	// Authors is the list of authors of the template repository.
	Authors []string `yaml:"authors"`

	// License is the license of the template repository.
	License string `yaml:"license"`

	// Variables is the list of variables of the template repository.
	Variables []CloneyVariable `yaml:"variables"`
}

// NewCloneyMetadataFromRawYAML creates a new CloneyMetadata struct from a YAML string.
func NewCloneyMetadataFromRawYAML(rawYAML string) (*CloneyMetadata, error) {
	// Parse YAML.
	var metadata CloneyMetadata
	err := yaml.Unmarshal([]byte(rawYAML), &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// GetVariablesMap returns a map of the variables in the Cloney template repository.
func (m *CloneyMetadata) GetVariablesMap() (map[string]interface{}, error) {
	// In a real implementation, you should return dynamic variables based on the metadata.
	// This example provides hardcoded dummy variables.
	variablesMap := make(map[string]interface{})
	variablesMap["app_name"] = "MyApp"
	variablesMap["enable_logging"] = true
	variablesMap["port"] = 8080
	return variablesMap, nil
}

// Show prints the Cloney template repository metadata in a pretty way.
func (m *CloneyMetadata) Show() {
	// Print general information table.
	fmt.Print("\nGeneral Information:\n\n")
	generalInfoTable := tablewriter.NewWriter(os.Stdout)
	generalInfoTable.SetHeader([]string{"Name", "Description", "Version", "Authors", "License"})
	generalInfoTable.Append(
		[]string{m.Name, m.Description, m.Version, strings.Join(m.Authors, ", "), m.License},
	)
	generalInfoTable.Render()

	// Print variables table.
	if len(m.Variables) == 0 {
		fmt.Print("\nNo variables provided.\n")
		return
	}
	fmt.Print("\nInput Variables:\n\n")
	variablesTable := tablewriter.NewWriter(os.Stdout)
	variablesTable.SetHeader([]string{"Name", "Description", "Type", "Default"})
	for _, variable := range m.Variables {
		variablesTable.Append(
			[]string{variable.Name, variable.Description, variable.Type, variable.Default},
		)
	}
	variablesTable.Render()
	fmt.Println()
}
