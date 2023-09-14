package metadata

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	tw "github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

// CloneyMetadataVariable represents a variable in a Cloney template repository.
type CloneyMetadataVariable struct {
	// Name is the name of the variable.
	Name string `yaml:"name" validate:"required"`

	// Description is the description of the variable.
	Description string `yaml:"description"`

	// Default is the default value of the variable.
	Default interface{} `yaml:"default"`

	// Example is an example value of the variable.
	Example interface{} `yaml:"example" validate:"required"`
}

// CloneyMetadata represents the metadata file of a Cloney template repository.
type CloneyMetadata struct {
	// Name is the name of the template repository.
	Name string `yaml:"name" validate:"required"`

	// Description is the description of the template repository.
	Description string `yaml:"description"`

	// Version is the version of the template repository.
	Version string `yaml:"version" validate:"required,semver"`

	// Authors is the list of authors of the template repository.
	Authors []string `yaml:"authors"`

	// License is the license of the template repository.
	License string `yaml:"license"`

	// Variables is the list of variables of the template repository.
	Variables []CloneyMetadataVariable `yaml:"variables"`
}

// NewCloneyMetadataFromRawYAML creates a new CloneyMetadata struct from a YAML string.
func NewCloneyMetadataFromRawYAML(rawYAML string) (*CloneyMetadata, error) {
	// Parse YAML.
	var metadata CloneyMetadata
	err := yaml.Unmarshal([]byte(rawYAML), &metadata)
	if err != nil {
		return nil, err
	}

	// Validate metadata.
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(metadata)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata file structure: %w", err)
	}

	// Validate variables separately because 'validator' package does not validate struct slices.
	for _, variable := range metadata.Variables {
		err = validate.Struct(variable)
		if err != nil {
			return nil, fmt.Errorf("invalid metadata file structure: %w", err)
		}
	}

	return &metadata, nil
}

// ValidateVariablesMatch validates if a given map of variables matches the variables defined
// in the template repository metadata file.
func (m *CloneyMetadata) ValidateVariablesMatch(userVariables map[string]interface{}) error {
	// Convert the template variables to a map.
	templateVariables := make(map[string]interface{})
	for _, variable := range m.Variables {
		// User variable types are compared to the type of the example
		// variables defined in the template metadata file.
		templateVariables[variable.Name] = variable.Example
	}

	// Check if the user variables match the template variables.
	if !reflect.DeepEqual(templateVariables, userVariables) {
		return fmt.Errorf("user variables do not match the template variables")
	}

	return nil
}

// ShowGeneralInformation prints the general information of the Cloney template repository in a pretty way.
func (m *CloneyMetadata) ShowGeneralInformation() {
	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Version", "Authors", "License"})
	table.SetHeaderColor(
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
	)
	table.SetAlignment(tw.ALIGN_LEFT)
	table.Append(
		[]string{
			m.Name,
			m.Description,
			m.Version,
			strings.Join(m.Authors, ", "),
			m.License,
		},
	)
	table.Render()
}

// ShowVariables prints the variables of the Cloney template repository in a pretty way.
func (m *CloneyMetadata) ShowVariables() {
	if len(m.Variables) == 0 {
		fmt.Println("This template repository has no variables.")
		return
	}
	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Type", "Default", "YAML Example"})
	table.SetHeaderColor(
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
	)
	table.SetAlignment(tw.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	for _, variable := range m.Variables {
		table.Append(
			[]string{
				variable.Name,
				variable.Description,
				VariableType(variable.Example),
				VariableValue(variable.Default),
				VariableValue(variable.Example),
			},
		)
	}
	table.Render()
}

// Show prints the Cloney template repository metadata in a pretty way.
func (m *CloneyMetadata) Show() {
	fmt.Println("General information about this template repository:")
	m.ShowGeneralInformation()
	fmt.Print("\nVariables of this template repository:\n")
	m.ShowVariables()
}
