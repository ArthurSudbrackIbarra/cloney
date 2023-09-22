package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"

	tw "github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
)

// testInfoCommand represents a command instance used for testing.
var testInfoCommand = CreateInfoCommand()

// createDummyCloneyMetadataFile creates a dummy Cloney metadata file in the specified directory.
func createDummyCloneyMetadataFile(assert *assert.Assertions, directory string) {
	// Define the raw Cloney metadata content.
	rawMetadata := `
manifest_version: v1
name: TestProject
description: This is a test project.
license: MIT
template_version: "0.0.0"
authors:
  - John Doe
variables:
  - name: test_variable
    description: This is a test variable.
    default: test
    example: test
`
	// Write the metadata content to a file in the specified directory.
	err := os.WriteFile(appConfig.MetadataFileName, []byte(rawMetadata), os.ModePerm)
	assert.NoError(err)
}

// createDummyGeneralInfoTable creates a dummy general Cloney project table in the provided buffer with mocked data.
func createDummyGeneralInfoTable(buffer bytes.Buffer) {
	table := tw.NewWriter(&buffer)

	// Set table headers with formatting.
	table.SetHeader([]string{"Name", "Description", "Template Version", "Authors", "License"})
	table.SetHeaderColor(
		tw.Colors{tw.Bold, tw.BgBlueColor},
		tw.Colors{tw.Bold, tw.BgBlueColor},
		tw.Colors{tw.Bold, tw.BgBlueColor},
		tw.Colors{tw.Bold, tw.BgBlueColor},
		tw.Colors{tw.Bold, tw.BgBlueColor},
	)
	table.SetAlignment(tw.ALIGN_LEFT)

	// Append row with dummy data.
	table.Append([]string{"TestProject", "This is a test project.", "0.0.0", "John Doe", "MIT"})

	// Render the table.
	table.Render()
}

// TestInfoCommandWhenCurrentDirectoryIsACloneyProject tests the output of the info command
// when the current directory is a Cloney project. It should not return an error.
func TestInfoCommandWhenCurrentDirectoryIsACloneyProject(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Create a dummy Cloney metadata file in the current directory.
	// This simulates a Cloney project.
	createDummyCloneyMetadataFile(assert, ".")

	// Redirect stdout to a buffer.
	var buffer bytes.Buffer
	testVersionCmd.SetOut(&buffer)

	// Execute the command.
	err := testInfoCommand.Execute()

	// Assert that the command did not return an error.
	assert.Nil(err)

	// Get the command output.
	cmdOutput := buffer.String()

	// Assert that the command printed the correct output.
	// To verify this, we create a table with the same data as the command output
	// and compare the two tables. The command output should contain the table output.
	buffer.Reset()
	createDummyGeneralInfoTable(buffer)
	mainTableOutput := buffer.String()
	assert.True(strings.Contains(mainTableOutput, cmdOutput))

	// Delete the dummy Cloney metadata file after the test.
	os.Remove(appConfig.MetadataFileName)
}
