package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testInfoCommand represents a command instance used for testing.
var testInfoCommand = CreateInfoCommand()

// CreateDummyCloneyMetadataFile creates a dummy Cloney metadata file in the specified directory.
func CreateDummyCloneyMetadataFile(assert *assert.Assertions, directory string) {
	// Define the raw Cloney metadata content.
	rawMetadata := `
manifest_version: v1
name: TestProject
description: This is a test project.
license: MIT
template_version: 0.0.0
authors:
  - John Doe
variables:
  - name: app_name
    description: The name of the application.
    example: My Bank App
  - name: dark_mode
    description: Whether or not to enable dark mode.
    example: false
  - name: currencies
    description: A list of currencies
    example:
      - name: USD
        symbol: abc
        description: United States Dollar	
`
	// Write the metadata content to a file in the specified directory.
	err := os.MkdirAll(directory, os.ModePerm)
	assert.NoError(err)
	metadataFilePath := filepath.Join(directory, appConfig.MetadataFileName)
	err = os.WriteFile(metadataFilePath, []byte(rawMetadata), os.ModePerm)
	assert.NoError(err)
}

// TestInfoCommandWhenCurrentDirectoryIsACloneyProject tests the "info" command
// when the current directory is a Cloney project. It should not return an error.
func TestInfoCommandWhenCurrentDirectoryIsACloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Create a dummy Cloney metadata file in the current directory to simulate a Cloney project.
	CreateDummyCloneyMetadataFile(assert, ".")

	// Execute the "info" command.
	err := testInfoCommand.Execute()

	// Assert that the "info" command did not return an error.
	assert.Nil(err)

	// Delete the dummy Cloney metadata file after the test.
	os.Remove(appConfig.MetadataFileName)
}

// TestInfoCommandWhenCurrentDirectoryIsNotACloneyProject tests the "info" command
// when the current directory is not a Cloney project. It should return an error.
func TestInfoCommandWhenCurrentDirectoryIsNotACloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Execute the "info" command.
	err := testInfoCommand.Execute()

	// Assert that the "info" command returned an error.
	assert.NotNil(err)
}

// TestInfoCommandPointingToLocalCloneyProject tests the output of the "info" command
// when the user specifies a local Cloney project directory.
// It should not return an error.
func TestInfoCommandPointingToLocalCloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Create a directory to simulate a Cloney project.
	err := os.Mkdir("test-project", os.ModePerm)
	assert.NoError(err)

	// Create a dummy Cloney metadata file in the test project directory.
	CreateDummyCloneyMetadataFile(assert, "test-project")

	// Simulate CLI arguments with flags and values to specify the project directory.
	testInfoCommand.SetArgs([]string{"./test-project"})

	// Execute the "info" command.
	err = testInfoCommand.Execute()

	// Assert that the "info" command did not return an error.
	assert.Nil(err)

	// Delete the created directory after the test.
	os.RemoveAll("test-project")
}

// TestInfoCommandPointingToRemoteCloneyProject tests the "info" command
// when the user specifies a valid remote Cloney project URL.
// It should not return an error.
func TestInfoCommandPointingToRemoteCloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Simulate CLI arguments with flags and values to specify the remote project URL.
	testInfoCommand.SetArgs([]string{appConfig.CloneyExampleRepositoryURL})

	// Execute the "info" command.
	err := testInfoCommand.Execute()

	// Assert that the "info" command did not return an error.
	assert.Nil(err)
}

// TestInfoCommandPointingToPrivateRemoteCloneyProject tests the "info" command
// when the user specifies a valid remote Cloney project URL that requires authentication.
// It should not return an error.
func TestInfoCommandPointingToPrivateRemoteCloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Simulate CLI arguments with flags and values to specify the remote project URL.
	// Pointing to a private repository using a personal access token.
	testInfoCommand.SetArgs([]string{
		"https://github.com/ArthurSudbrackIbarra/private-cloney.git",
		"--token", os.Getenv("PERSONAL_ACCESS_TOKEN"),
	})

	// Execute the "info" command.
	err := testInfoCommand.Execute()

	// Assert that the "info" command did not return an error.
	assert.Nil(err)

	// Reset the command flags.
	ResetInfoCommandFlags(testInfoCommand)
}

// TestInfoCommandPointingToInvalidRemoteCloneyProject tests the "info" command
// when the user specifies an invalid remote Cloney project URL.
// It should return an error.
func TestInfoCommandPointingToInvalidRemoteCloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Simulate CLI arguments with flags and values to specify an invalid remote project URL.
	testInfoCommand.SetArgs([]string{"https://github.com/ArthurSudbrackIbarra/cloney.git"})

	// Execute the "info" command.
	err := testInfoCommand.Execute()

	// Assert that the "info" command returned an error.
	assert.NotNil(err)
}
