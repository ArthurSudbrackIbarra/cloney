package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// NewCloneyUserVariablesFromFile reads a file and returns a map with the variables defined in it.
// It accepts files with the following extensions: '.json', '.yaml' or '.yml'.
func NewCloneyUserVariablesFromFile(filePath string) (map[string]interface{}, error) {
	// Read file content.
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read variables file: %w", err)
	}

	variables := make(map[string]interface{})

	// If the file extension is '.json', parse it as JSON.
	// If the file extension is '.yaml' or '.yml', parse it as YAML.
	// Otherwise, return an error.
	if strings.HasSuffix(filePath, ".json") {
		err = json.Unmarshal(content, &variables)
		if err != nil {
			return nil, fmt.Errorf("failed to parse variables JSON file: %w", err)
		}
	} else if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		err = yaml.Unmarshal(content, &variables)
		if err != nil {
			return nil, fmt.Errorf("failed to parse variables YAML file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file extension, expected '.json', '.yaml' or '.yml'")
	}

	return variables, nil
}

// NewCloneyUserVariablesFromRawJSON returns a map with the variables defined in the given JSON string.
func NewCloneyUserVariablesFromRawJSON(rawJSON string) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	err := json.Unmarshal([]byte(rawJSON), &variables)
	if err != nil {
		return nil, fmt.Errorf("failed to parse variables JSON string: %w", err)
	}

	return variables, nil
}
