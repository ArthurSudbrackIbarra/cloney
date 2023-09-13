package metadata

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kr/pretty"
)

// Constants for variable types.
const STRING_VARIABLE_TYPE = "string"
const INTEGER_VARIABLE_TYPE = "integer"
const DECIMAL_VARIABLE_TYPE = "decimal"
const BOOLEAN_VARIABLE_TYPE = "boolean"
const LIST_VARIABLE_TYPE = "list"
const MAP_VARIABLE_TYPE = "map"
const UNKNOWN_VARIABLE_TYPE = "unknown"

// MapVariableToString converts a variable of type map to a string representation.
func MapVariableToString(value interface{}) string {
	// Define regular expressions for substitutions.

	// This regex transforms '[]interface' {} into 'list'.
	regexList := regexp.MustCompile(`\[\]interface \{\}`)
	// This regex transforms 'map[interface {}]interface {}' into 'map'.
	regexMap := regexp.MustCompile(`map\[interface \{\}\]interface \{\}`)

	// This regex transforms '"key": "value"' into 'group + (string)'.
	regexString := regexp.MustCompile(`(".+":) ".+"`)
	// This regex transforms '"key": [0-9]+' into 'group + (integer)'.
	regexInteger := regexp.MustCompile(`(".+":) [0-9]+`)
	// This regex transforms '"key": [0-9]+\.[0-9]+' into 'group + (decimal)'.
	regexDecimal := regexp.MustCompile(`(".+":) [0-9]+\.[0-9]+`)
	// This regex transforms '"key": true|false' into 'group + (boolean)'.
	regexBoolean := regexp.MustCompile(`(".+":) \btrue|false\b`)

	// This regex transforms 'string(value)' into 'string'.
	regexString2 := regexp.MustCompile(`string\(.+\)`)
	// This regex transforms 'int(value)' into 'integer'.
	regexInteger2 := regexp.MustCompile(`int\(.+\)`)
	// This regex transforms 'float(value)' into 'decimal'.
	regexDecimal2 := regexp.MustCompile(`float\(.+\)`)
	// This regex transforms 'bool(value)' into 'boolean'.
	regexBoolean2 := regexp.MustCompile(`bool\(.+\)`)

	// Use 'pretty' to get a string representation of the variable.
	dumped := pretty.Sprint(value)

	// Remove commas.
	dumped = strings.ReplaceAll(dumped, ",", "")

	// Apply the substitutions.
	dumped = regexList.ReplaceAllString(dumped, fmt.Sprintf("%s ", LIST_VARIABLE_TYPE))
	dumped = regexMap.ReplaceAllString(dumped, fmt.Sprintf("%s ", MAP_VARIABLE_TYPE))
	dumped = regexString.ReplaceAllString(dumped, fmt.Sprintf("$1 %s", STRING_VARIABLE_TYPE))
	dumped = regexInteger.ReplaceAllString(dumped, fmt.Sprintf("$1 %s", INTEGER_VARIABLE_TYPE))
	dumped = regexDecimal.ReplaceAllString(dumped, fmt.Sprintf("$1 %s", DECIMAL_VARIABLE_TYPE))
	dumped = regexBoolean.ReplaceAllString(dumped, fmt.Sprintf("$1 %s", BOOLEAN_VARIABLE_TYPE))
	dumped = regexString2.ReplaceAllString(dumped, STRING_VARIABLE_TYPE)
	dumped = regexInteger2.ReplaceAllString(dumped, INTEGER_VARIABLE_TYPE)
	dumped = regexDecimal2.ReplaceAllString(dumped, DECIMAL_VARIABLE_TYPE)
	dumped = regexBoolean2.ReplaceAllString(dumped, BOOLEAN_VARIABLE_TYPE)

	return dumped
}
