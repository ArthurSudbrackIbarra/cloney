package metadata

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// NewCloneyUserVariablesFromFile reads a file and returns a map with the variables defined in it.
// It accepts files with the following extensions: '.yaml' or '.yml'.
func NewCloneyUserVariablesFromFile(filePath string) (map[string]interface{}, error) {
	// Read file content.
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read variables file: %w", err)
	}

	variables := make(map[string]interface{})

	// If the file extension is '.yaml' or '.yml', parse it as YAML.
	// Otherwise, return an error.
	if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		err = yaml.Unmarshal(content, &variables)
		if err != nil {
			return nil, fmt.Errorf("failed to parse variables YAML file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file extension, expected '.yaml' or '.yml'")
	}

	return variables, nil
}

// NewCloneyUserVariablesFromRawYAML returns a map with the variables defined in the given raw YAML string.
func NewCloneyUserVariablesFromRawYAML(rawYAML string) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	err := yaml.Unmarshal([]byte(rawYAML), &variables)
	if err != nil {
		return nil, fmt.Errorf("failed to parse variables inline YAML: %w", err)
	}

	return variables, nil
}
