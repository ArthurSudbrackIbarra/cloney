package metadata

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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
	// mapVar can either be a map[string]interface{} or a map[interface{}]interface{}.
	// TODO: So we need to check which one it is.
	map_ := mapVar.(map[string]interface{})

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

// AreVariablesSameType checks if two variables are of the same type.
func AreVariablesSameType(variable1 interface{}, variable2 interface{}) bool {
	type1 := reflect.TypeOf(variable1)
	type2 := reflect.TypeOf(variable2)

	// First check if the types are the same.
	if type1 != type2 {
		return false
	}

	// If the types are the same, check if they are slices.
	if type1.Kind() == reflect.Slice {
		if type1.Elem() != type2.Elem() {
			return false
		}

		value1 := reflect.ValueOf(variable1)
		value2 := reflect.ValueOf(variable2)

		// If the slices are empty, return true.
		if value1.Len() == 0 && value2.Len() == 0 {
			return true
		}

		// Check if the elements at position 0 are of the same type.
		element1 := value1.Index(0)
		element2 := value2.Index(0)
		return AreVariablesSameType(element1.Interface(), element2.Interface())
	}

	// If the types are maps, check if the keys are the same type and values are of the same type recursively.
	if type1.Kind() == reflect.Map {
		keyType1 := type1.Key()
		keyType2 := type2.Key()
		valueType1 := type1.Elem()
		valueType2 := type2.Elem()

		if keyType1 != keyType2 || valueType1 != valueType2 {
			return false
		}

		mapValue1 := reflect.ValueOf(variable1)
		mapValue2 := reflect.ValueOf(variable2)

		// Check if the keys are the same in both maps.
		for _, key := range mapValue1.MapKeys() {
			if !mapValue2.MapIndex(key).IsValid() {
				return false
			}
		}

		// Check if the values are of the same type for each key.
		for _, key := range mapValue1.MapKeys() {
			if !AreVariablesSameType(mapValue1.MapIndex(key).Interface(), mapValue2.MapIndex(key).Interface()) {
				return false
			}
		}

		return true
	}

	// If the types are not slices or maps, they are basic types (ints, floats, bools, strings).
	return true
}
