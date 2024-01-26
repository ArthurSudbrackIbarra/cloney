package metadata

import (
	"fmt"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/pkg/terminal"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// CloneyMetadataConfiguration represents the configuration of a Cloney template repository.
type CloneyMetadataConfiguration struct {
	// IgnorePaths is a list of paths to ignore when cloning the template repository.
	IgnorePaths []string `yaml:"ignore_paths"`

	// PostCloneCommands is a list of commands to run after cloning the template repository.
	// Each command is a list of strings, where the first element is the command name and the rest are the arguments.
	// Example: [["echo", "Hello World!"]]
	PostCloneCommands [][]string `yaml:"post_clone_commands"`

	// UnknownFields captures all unknown fields within configuration and must be empty after parsing.
	UnknownFields map[string]interface{} `yaml:",inline"`
}

// CloneyMetadataVariable represents a variable in a Cloney template repository.
type CloneyMetadataVariable struct {
	// Name is the variable name.
	Name string `yaml:"name" validate:"required"`

	// Description is the variable description.
	Description string `yaml:"description"`

	// Default is the default value of the variable.
	Default interface{} `yaml:"default"`

	// Example is an example value of the variable.
	Example interface{} `yaml:"example" validate:"required"`

	// Validate specifies if the variable should be validated.
	// It is a pointer to a bool because if the field is not defined in the YAML file,
	// the default value should be true.
	Validate *bool `yaml:"validate"`

	// UnknownFields captures all unknown fields within variables and must be empty after parsing.
	UnknownFields map[string]interface{} `yaml:",inline"`
}

// CloneyMetadata represents the metadata file of a Cloney template repository.
type CloneyMetadata struct {
	// Name is the template repository name.
	Name string `yaml:"name" validate:"required"`

	// Description is the template repository description.
	Description string `yaml:"description"`

	// TemplateVersion is the version of the template repository.
	TemplateVersion string `yaml:"template_version" validate:"required,semver"`

	// ManifestVersion is the version of the manifest file.
	ManifestVersion string `yaml:"manifest_version" validate:"required"`

	// Authors is the list of authors of the template repository.
	Authors []string `yaml:"authors"`

	// License is the license of the template repository.
	License string `yaml:"license"`

	// Configuration is the configuration of the template repository.
	Configuration CloneyMetadataConfiguration `yaml:"configuration"`

	// Variables is the list of variables of the template repository.
	Variables []CloneyMetadataVariable `yaml:"variables"`

	// UnknownFields captures all unknown fields at root level and must be empty after parsing.
	UnknownFields map[string]interface{} `yaml:",inline"`
}

// NewCloneyMetadataFromRawYAML creates a new CloneyMetadata struct from a YAML string.
// It also validates the manifest version and the metadata structure.
func NewCloneyMetadataFromRawYAML(rawYAML string, supportedManifestVersions []string) (*CloneyMetadata, []error) {
	// List to store errors found during validation.
	// This is to return all errors at once instead of returning the first error found.
	errors := make([]error, 0)

	// Parse YAML.
	var metadata CloneyMetadata
	err := yaml.Unmarshal([]byte(rawYAML), &metadata)
	if err != nil {
		errors = append(errors, fmt.Errorf("Invalid YAML: %w", err))
		return nil, errors
	}

	// Check for unknown fields at the root level.
	if len(metadata.UnknownFields) > 0 {
		errors = append(errors, fmt.Errorf("Unknown field(s) at root level: %s", getKeys(metadata.UnknownFields)))
	}

	// Validate metadata.
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(metadata)
	if err != nil {
		// Custom error message for some validation errors.
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			// Convert field name to lowercase.
			// If the field name is 'manifestversion', convert it to 'manifest_version'.
			// If the field name is 'templateversion', convert it to 'template_version'.
			var field string
			field = strings.ToLower(validationError.Field())
			if field == "manifestversion" {
				field = "manifest_version"
			} else if field == "templateversion" {
				field = "template_version"
			}
			switch validationError.Tag() {
			case "required":
				errors = append(errors, fmt.Errorf("Missing required field '%s' at root level", field))
			case "semver":
				errors = append(errors, fmt.Errorf("Invalid semantic version '%s' for field %s", validationError.Value(), field))
			default:
				{
					errors = append(errors, fmt.Errorf("Invalid metadata file structure: %w", err))
					return nil, errors
				}
			}
		}
	}

	// Check if manifest version is supported.
	versionSupported := false
	for _, supportedManifestVersion := range supportedManifestVersions {
		if metadata.ManifestVersion == supportedManifestVersion {
			versionSupported = true
			break
		}
	}
	if !versionSupported {
		errors = append(errors, fmt.Errorf(
			"Manifest version '%s' is not supported in this Cloney version.\nPlease update or downgrade your Cloney version.\n\nSupported versions: %s",
			metadata.ManifestVersion,
			strings.Join(supportedManifestVersions, ", "),
		))
		return nil, errors
	}

	// Validate the configuration block.
	// If manifest_version is v1, post_clone_commands field is not supported.
	if metadata.ManifestVersion == "v1" && len(metadata.Configuration.PostCloneCommands) > 0 {
		errors = append(errors, fmt.Errorf("Manifest version '%s' does not support the 'configuration.post_clone_commands' field, please update the manifest version to 'v2'", metadata.ManifestVersion))
	}

	// Check for unknown fields in Configuration.
	if len(metadata.Configuration.UnknownFields) > 0 {
		errors = append(errors, fmt.Errorf("Unknown field(s) in the configuration block: %s", getKeys(metadata.Configuration.UnknownFields)))
	}

	// Validate variables separately because 'validator' package does not validate struct slices.
	for _, variable := range metadata.Variables {
		err = validate.Struct(variable)
		if err != nil {
			// Custom error message for some validation errors.
			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				switch validationError.Tag() {
				case "required":
					errors = append(errors, fmt.Errorf(
						"Missing required field '%s' for variable '%s'",
						strings.ToLower(validationError.Field()),
						variable.Name,
					))
				default:
					{
						errors = append(errors, fmt.Errorf("Invalid variable '%s' structure: %w", variable.Name, err))
					}
				}
			}
		}

		// If the variable has a default value, check if it is of the same type as the example value.
		if variable.Default != nil && variable.Example != nil && !AreVariablesSameType(variable.Example, variable.Default) {
			errors = append(errors, fmt.Errorf(
				"Variable '%s' has a default value of type '%s' but its example value is of type '%s'",
				variable.Name,
				VariableType(variable.Default),
				VariableType(variable.Example),
			))
		}

		// Check for unknown fields in Variables.
		if len(variable.UnknownFields) > 0 {
			errors = append(errors, fmt.Errorf("Unknown field(s) for variable '%s': %s", variable.Name, getKeys(variable.UnknownFields)))
		}
	}

	// If there are errors, return them.
	if len(errors) > 0 {
		return nil, errors
	}

	return &metadata, nil
}

// getKeysString returns the keys of a map as a string with comma-separated values surrounded by single quotes.
func getKeys(m map[string]interface{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, fmt.Sprintf("'%s'", k))
	}
	return strings.Join(keys, ", ")
}

// MatchUserVariables validates if a given map of variables matches the variables defined
// in the template repository's metadata file.
// It also adds the default values of the variables to the user variables if they are not defined.
func (m *CloneyMetadata) MatchUserVariables(userVariables map[string]interface{}) (map[string]interface{}, error) {
	// Iterate over the variables defined in the template repository's metadata file.
	for _, variable := range m.Variables {
		// Check if the variable is defined in the user variables.
		if _, contains := userVariables[variable.Name]; !contains {
			// If not defined, it could be that the variable is optional, so check if it has a default value.
			if variable.Default == nil {
				return nil, fmt.Errorf("variable '%s' is required but is not defined", variable.Name)
			} else {
				// If the variable has a default value, add it to the user variables.
				userVariables[variable.Name] = variable.Default
			}
		}

		// If the user specified that the variable should not be validated, skip validation.
		if variable.Validate != nil && !*variable.Validate {
			continue
		}

		// Check if the variables are of the same type.
		// User variables are compared against the example value of the template repository's metadata file.
		type1 := VariableType(variable.Example)
		type2 := VariableType(userVariables[variable.Name])

		// Special case, if the template variable is a 'decimal' and the user variable is an 'integer', it is valid.
		// This is because integers are a subset of decimals.
		if type1 == DECIMAL_VARIABLE_TYPE && type2 == INTEGER_VARIABLE_TYPE {
			continue
		}

		// In other cases, if the types are different, return an error.
		if !AreVariablesSameType(variable.Example, userVariables[variable.Name]) {
			return nil, fmt.Errorf("variable '%s' is of type '%s' but must be of type '%s'", variable.Name, type2, type1)
		}
	}
	return userVariables, nil
}

// GetGeneralInfo returns the general information of the Cloney template repository as a string.
func (m *CloneyMetadata) GetGeneralInfo() string {
	result := terminal.WhiteBoldUnderline("\nGeneral Information\n\n")
	result += fmt.Sprintf("%s: %s\n", "Template Name", m.Name)
	result += fmt.Sprintf("%s: %s\n", "Template Description", m.Description)
	result += fmt.Sprintf("%s: %s\n", "Template Version", m.TemplateVersion)
	result += fmt.Sprintf("%s: %s\n", "Template License", m.License)
	result += fmt.Sprintf("%s: %s\n", "Template Author(s)", strings.Join(m.Authors, ", "))
	return result
}

// GetVariables returns the variables of the Cloney template repository as a string.
func (m *CloneyMetadata) GetVariables() string {
	result := "\n"
	for index, variable := range m.Variables {
		if variable.Default == nil {
			result += fmt.Sprintf("%s %s\n\n", terminal.WhiteBoldUnderline("Variable"), fmt.Sprintf("%s (%s)", terminal.BlueBoldUnderline(variable.Name), terminal.Yellow("Required")))
		} else {
			result += fmt.Sprintf("%s %s (Optional)\n\n", terminal.WhiteBoldUnderline("Variable"), terminal.BlueBoldUnderline(variable.Name))
		}

		result += fmt.Sprintf("%s: %s\n", "Variable Description", variable.Description)

		varType := VariableType(variable.Example)
		if !strings.Contains(varType, "\n") {
			result += fmt.Sprintf("%s: %s\n", "Variable Type", VariableType(variable.Example))
		} else {
			result += fmt.Sprintf("%s:\n%s\n", "Variable Type", VariableType(variable.Example))
		}

		if variable.Default != nil {
			varDefault := VariableValue(variable.Default)
			if !strings.Contains(varDefault, "\n") {
				result += fmt.Sprintf("%s: %s\n", "Default Value", VariableValue(variable.Default))
			} else {
				result += fmt.Sprintf("%s:\n%s\n", "Default Value", VariableValue(variable.Default))
			}
		}

		varExample := VariableValue(variable.Example)
		if !strings.Contains(varExample, "\n") {
			result += fmt.Sprintf("%s: %s\n", "Example Value", VariableValue(variable.Example))
		} else {
			result += fmt.Sprintf("%s:\n%s\n", "Example Value", VariableValue(variable.Example))
		}
		if index != len(m.Variables)-1 {
			result += "\n"
		}
	}
	return result
}

// String returns the string representation of the CloneyMetadata struct.
func (m *CloneyMetadata) String() string {
	result := m.GetGeneralInfo()
	result += m.GetVariables()
	return result
}
