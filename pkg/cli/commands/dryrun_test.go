package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testDryRunCommand represents a command instance used for testing.
var testDryRunCommand = CreateDryRunCommand()

// CreateDummyVariablesFile creates a dummy Cloney variables file in the specified directory.
func CreateDummyVariablesFile(assert *assert.Assertions, directory string) {
	// Define the raw Cloney variables content.
	rawVariables := `
app_name: TestProject
dark_mode: true
currencies:
  - name: USD
    symbol: $
    description: "United States Dollar"
  - name: EUR
    symbol: €
    description: "Euro"
  - name: GBP
    symbol: £
    description: "British Pound"
`
	// Write the variables content to a file in the specified directory.
	err := os.MkdirAll(directory, os.ModePerm)
	assert.NoError(err)
	variablesFilePath := filepath.Join(directory, appConfig.DefaultUserVariablesFileName)
	err = os.WriteFile(variablesFilePath, []byte(rawVariables), os.ModePerm)
	assert.NoError(err)
}

// CreateDummyTxtFile creates a dummy txt file in the specified directory.
// Referencing the variables defined in CreateDummyVariablesFile.
func CreateDummyTxtFile(assert *assert.Assertions, directory string) {
	// Define the dummy file content
	rawTxt := `
This is a dummy file.
Your application name is {{ .app_name }}.
Dark mode is {{ .dark_mode }}.
Currencies:
{{- range .currencies }}
	- {{ .name }} ({{ .symbol }}): {{ .description }}
{{- end }}
`
	// Write the dummy file content to a file in the specified directory.
	err := os.MkdirAll(directory, os.ModePerm)
	assert.NoError(err)
	txtFilePath := filepath.Join(directory, "dummy.txt")
	err = os.WriteFile(txtFilePath, []byte(rawTxt), os.ModePerm)
	assert.NoError(err)
}

// TestDryRunCommandWhenCurrentDirectoryIsACloneyProject tests the "dry-run" command
// when the current directory is a Cloney project. It should not return an error.
func TestDryRunCommandWhenCurrentDirectoryIsACloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Create a dummy Cloney metadata file in the current directory.
	CreateDummyCloneyMetadataFile(assert, ".")

	// Create a dummy Cloney variables file in the current directory.
	CreateDummyVariablesFile(assert, ".")

	// Create a dummy txt file in the current directory.
	CreateDummyTxtFile(assert, ".")

	// Run the "dryrun" command.
	err := testDryRunCommand.Execute()

	// Assert that the "dryrun" command did not return an error.
	assert.Nil(err)

	// Assert that the dry run directory was created.
	assert.DirExists(appConfig.DefaultDryRunDirectoryName)

	// Assert that the cloney metadata file was not copied to the dry run directory.
	assert.NoFileExists(filepath.Join(appConfig.DefaultDryRunDirectoryName, appConfig.MetadataFileName))

	// Assert that the cloney variables file was not copied to the dry run directory.
	assert.NoFileExists(filepath.Join(appConfig.DefaultDryRunDirectoryName, appConfig.DefaultUserVariablesFileName))

	// Assert that the txt file was copied to the dry run directory.
	assert.FileExists(filepath.Join(appConfig.DefaultDryRunDirectoryName, "dummy.txt"))

	// Delete the dummy Cloney metadata file after the test.
	err = os.Remove(appConfig.MetadataFileName)
	assert.NoError(err)

	// Delete the dummy Cloney variables file after the test.
	err = os.Remove(appConfig.DefaultUserVariablesFileName)
	assert.NoError(err)

	// Delete the dummy txt file after the test.
	err = os.Remove("dummy.txt")
	assert.NoError(err)

	// Delete the dry run directory after the test.
	err = os.RemoveAll(appConfig.DefaultDryRunDirectoryName)
	assert.NoError(err)
}

// TODO: Add remaining tests...
