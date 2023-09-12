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

// GetType returns the type of a variable as a string.
func GetVariableType(variable interface{}) string {
	// The variable type is determined by checking the type of the example value.
	variableType := reflect.TypeOf(variable).String()

	// All types of integers are classified as "integer".
	if strings.HasPrefix(variableType, "int") {
		return "integer"
	}
	// All types of floats are classified as "decimal".
	if strings.HasPrefix(variableType, "float") {
		return "decimal"
	}
	// All booleans are classified as "boolean".
	if variableType == "bool" {
		return "boolean"
	}
	// All lists are classified as "list".
	if strings.HasPrefix(variableType, "[]") {
		return "list"
	}
	// All complex structures are classified as "map".
	if strings.HasPrefix(variableType, "map") {
		return "map"
	}

	// Otherwise, return the type as is.
	return variableType
}

// ValueToString returns the value of a variable as a string.
func ValueToString(value interface{}) string {
	// If the default value is nil, return an empty string.
	if value == nil {
		return ""
	}

	// Get the type of the variable.
	variableType := GetVariableType(value)

	// If value is a number or a boolean, return it as a string.
	if variableType == "integer" ||
		variableType == "decimal" ||
		variableType == "boolean" {
		return fmt.Sprintf("%v", value)
	}

	// If the value is a string, add quotes to it.
	if variableType == "string" {
		return fmt.Sprintf("\"%v\"", value)
	}

	// If the value is a list (slice), iterate over it and convert each element.
	if variableType == "list" {
		sliceValue := reflect.ValueOf(value)
		result := "["
		for i := 0; i < sliceValue.Len(); i++ {
			elementValue := sliceValue.Index(i).Interface()
			elementStr := ValueToString(elementValue)
			result += elementStr
			if i < sliceValue.Len()-1 {
				result += ", "
			}
		}
		result += "]"
		return result
	}

	// If the value is a map, return it as a string representation.
	if variableType == "map" {
		return fmt.Sprintf("%v", value)
	}

	return ""
}

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

	return &metadata, nil
}

// ==========================================

// This function will be REMOVED LATER...
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

// ==========================================

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
		[]string{m.Name, m.Description, m.Version, strings.Join(m.Authors, ", "), m.License},
	)
	table.Render()
}

// ShowVariables prints the variables of the Cloney template repository in a pretty way.
func (m *CloneyMetadata) ShowVariables() {
	if len(m.Variables) == 0 {
		fmt.Print("No variables provided.")
		return
	}
	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Type", "Default", "Example"})
	table.SetHeaderColor(
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
		tw.Colors{tw.Bold, tw.BgYellowColor, tw.FgBlackColor},
	)
	table.SetAlignment(tw.ALIGN_LEFT)
	for _, variable := range m.Variables {
		table.Append(
			[]string{variable.Name, variable.Description, GetVariableType(variable.Example), ValueToString(variable.Default), "TODO"},
		)
	}
	table.Render()
}

// Show prints the Cloney template repository metadata in a pretty way.
func (m *CloneyMetadata) Show() {
	m.ShowGeneralInformation()
	m.ShowVariables()
}
