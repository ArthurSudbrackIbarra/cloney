package metadata

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// Constants for variable types.
const (
	STRING_VARIABLE_TYPE  = "string"
	INTEGER_VARIABLE_TYPE = "integer"
	DECIMAL_VARIABLE_TYPE = "decimal"
	BOOLEAN_VARIABLE_TYPE = "boolean"
	LIST_VARIABLE_TYPE    = "list"
	MAP_VARIABLE_TYPE     = "map"
	UNKNOWN_VARIABLE_TYPE = "unknown"
)

// VariableType returns the type of a variable as a string.
func VariableType(variable interface{}) string {
	// Determine the variable type by checking the type of the example value.
	variableType := reflect.TypeOf(variable).String()

	// Check for common types and classify them.
	switch variableType {
	case "string":
		return STRING_VARIABLE_TYPE
	case "int", "int8", "int16", "int32", "int64":
		return INTEGER_VARIABLE_TYPE
	case "float32", "float64":
		return DECIMAL_VARIABLE_TYPE
	case "bool":
		return BOOLEAN_VARIABLE_TYPE
	}

	// Handle maps and lists separately.
	if strings.HasPrefix(variableType, "map") {
		return MapVariableType(variable)
	} else if strings.HasPrefix(variableType, "[]") {
		return ListVariableType(variable)
	}

	// Return "unknown" for other types.
	return UNKNOWN_VARIABLE_TYPE
}

// VariableValue returns the value of a variable as a string.
func VariableValue(value interface{}) string {
	// If the value is nil, return an empty string.
	if value == nil {
		return ""
	}

	// Get the type of the variable.
	variableType := VariableType(value)

	// If value is a number, a boolean or a string, convert it to a string.
	if variableType == INTEGER_VARIABLE_TYPE ||
		variableType == DECIMAL_VARIABLE_TYPE ||
		variableType == BOOLEAN_VARIABLE_TYPE ||
		variableType == STRING_VARIABLE_TYPE {
		return fmt.Sprintf("%v", value)
	}

	// If the variable is not of the types above, and it is not unknown, then it is either a map or a list.
	// In this case, we convert it to a YAML string.
	if variableType != UNKNOWN_VARIABLE_TYPE {
		valueYAML, _ := yaml.Marshal(value)
		return string(valueYAML)
	}

	// Otherwise, return "".
	return ""
}

// ListVariableType returns a string representation of the type of a list variable.
func ListVariableType(listVar interface{}) string {
	// Check if it's a list.
	variableType := reflect.TypeOf(listVar).String()
	if !strings.HasPrefix(variableType, "[]") {
		return UNKNOWN_VARIABLE_TYPE
	}

	// Convert to a Go slice.
	slice := listVar.([]interface{})

	if len(slice) == 0 {
		return LIST_VARIABLE_TYPE
	}

	// Get the type of the first value in the list.
	firstValueType := VariableType(slice[0])

	// Assuming all values in the list have the same type, return the type idented.
	return IndentStringStructure(
		fmt.Sprintf("%s [\n%s\n]", LIST_VARIABLE_TYPE, firstValueType),
	)
}

// MapVariableType returns a string representation of the type of a map variable.
func MapVariableType(mapVar interface{}) string {
	// Check if it's a map.
	variableType := reflect.TypeOf(mapVar).String()
	if !strings.HasPrefix(variableType, "map") {
		return UNKNOWN_VARIABLE_TYPE
	}

	// Convert to a Go map.
	map_ := mapVar.(map[interface{}]interface{})

	if len(map_) == 0 {
		return MAP_VARIABLE_TYPE
	}

	// Iterate over the map values and get their types.
	variableTypes := ""
	for key, value := range map_ {
		valueType := VariableType(value)

		// If any value has an unknown type, return "unknown" for the map.
		if valueType == UNKNOWN_VARIABLE_TYPE {
			return UNKNOWN_VARIABLE_TYPE
		}

		variableTypes += fmt.Sprintf("%s: %s\n", key, valueType)
	}

	// Return the map type indented.
	return IndentStringStructure(
		fmt.Sprintf("%s {\n%s\n}", MAP_VARIABLE_TYPE, variableTypes),
	)
}

// IndentStringStructure indents the structure of a string (list or map-like structure).
// Used for formatting maps and lists.
func IndentStringStructure(input string) string {
	lines := strings.Split(input, "\n")
	indentLevel := 0
	output := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, "{") || strings.Contains(line, "[") {
			output += strings.Repeat("  ", indentLevel) + line + "\n"
			indentLevel++
		} else if strings.Contains(line, "}") || strings.Contains(line, "]") {
			indentLevel--
			output += strings.Repeat("  ", indentLevel) + line + "\n"
		} else {
			output += strings.Repeat("  ", indentLevel) + line + "\n"
		}
	}

	return output
}
