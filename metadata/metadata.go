package metadata

import (
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
