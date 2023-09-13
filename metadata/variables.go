package metadata

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// Constants for variable types.
const STRING_VARIABLE_TYPE = "string"
const INTEGER_VARIABLE_TYPE = "integer"
const DECIMAL_VARIABLE_TYPE = "decimal"
const BOOLEAN_VARIABLE_TYPE = "boolean"
const LIST_VARIABLE_TYPE = "list"
const MAP_VARIABLE_TYPE = "map"
const UNKNOWN_VARIABLE_TYPE = "unknown"

// VariableType returns the type of a variable as a string.
func VariableType(variable interface{}) string {
	// The variable type is determined by checking the type of the example value.
	variableType := reflect.TypeOf(variable).String()

	// All types of integers are classified as "integer".
	if strings.HasPrefix(variableType, "int") {
		return INTEGER_VARIABLE_TYPE
	}
	// All types of floats are classified as "decimal".
	if strings.HasPrefix(variableType, "float") {
		return DECIMAL_VARIABLE_TYPE
	}
	// All booleans are classified as "boolean".
	if variableType == "bool" {
		return BOOLEAN_VARIABLE_TYPE
	}
	// All strings are classified as "string".
	if variableType == "string" {
		return STRING_VARIABLE_TYPE
	}
	// Maps have their own way of being classified.
	// We call MapVariableToString to get a string representation of the map.
	if strings.HasPrefix(variableType, "map") {
		return MapVariableType(variable)
	}

	// Lists also have their own way of being classified.
	// We call ListVariableToString to get a string representation of the list.
	if strings.HasPrefix(variableType, "[]") {
		return ListVariableType(variable)
	}

	// Otherwise, return "unknown".
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
	// If not a list, return "unknown".
	variableType := reflect.TypeOf(listVar).String()
	if !strings.HasPrefix(variableType, "[]") {
		return UNKNOWN_VARIABLE_TYPE
	}

	// Convert to a Go slice.
	slice := listVar.([]interface{})

	// If the list is empty, return just "list".
	// But this case will never happen in a real template repository
	// Because the variable[index].example will never be an empty list.
	if len(slice) == 0 {
		return LIST_VARIABLE_TYPE
	}

	// Get the type of the first value.
	firstValueType := VariableType(slice[0])

	// Assuming all values of the list have the same type,
	// return the type of the first value.
	return IndentStringStructure(
		fmt.Sprintf("%s [\n%s\n]", LIST_VARIABLE_TYPE, firstValueType),
	)
}

// MapVariableType returns a string representation of the type of a map variable.
func MapVariableType(mapVar interface{}) string {
	// If not a map, return "unknown".
	variableType := reflect.TypeOf(mapVar).String()
	if !strings.HasPrefix(variableType, "map") {
		return UNKNOWN_VARIABLE_TYPE
	}

	// Convert to a Go map.
	// Name it 'map_' because 'map' is a reserved word.
	map_ := mapVar.(map[interface{}]interface{})

	// If the map is empty, return just "map".
	// But this case will never happen in a real template repository
	// Because the variable.key.example will never be an empty map.
	if len(map_) == 0 {
		return MAP_VARIABLE_TYPE
	}

	// Iterate over the map values and get their types.
	variableTypes := ""
	for key, value := range map_ {
		// Get the type of the value.
		valueType := VariableType(value)

		// If the type is unknown, return "unknown".
		if valueType == UNKNOWN_VARIABLE_TYPE {
			return UNKNOWN_VARIABLE_TYPE
		}

		// Add the type to the map.
		variableTypes += fmt.Sprintf("%s: %s\n", key, valueType)
	}

	return IndentStringStructure(
		fmt.Sprintf("%s {\n%s\n}", MAP_VARIABLE_TYPE, variableTypes),
	)
}

// IndentStringStructure indents the structure of a string.
// Used for identing maps and lists.
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
