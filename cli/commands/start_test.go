package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// testStartCmd represents a command instance used for testing.
var testStartCmd = CreateStartCommand()

// TestMetadataProperties is a struct to store the metadata properties.
type TestMetadataProperties struct {
	Name            string
	Description     string
	License         string
	TemplateVersion string
	ManifestVersion string
}

// assertMetadataMainFieldsAreOK asserts that the metadata file contains the main fields
// and that the fields have the given values.
func assertMetadataMainFieldsAreOK(assert *assert.Assertions, metadataProperties TestMetadataProperties, parsedMetadata map[string]interface{}) {
	assert.Equal(metadataProperties.Name, parsedMetadata["name"])
	assert.Equal(metadataProperties.Description, parsedMetadata["description"])
	assert.Equal(metadataProperties.License, parsedMetadata["license"])
	assert.Equal(metadataProperties.TemplateVersion, parsedMetadata["template_version"])
	assert.Equal(metadataProperties.ManifestVersion, parsedMetadata["manifest_version"])
}

// assertVariablesAreOK asserts that the metadata file contains the example variables
// and that the variables have the correct values.
func assertMetadataVariablesAreOK(assert *assert.Assertions, parsedMetadata map[string]interface{}) {
	// Assert that the metadata file contains the variables.
	assert.Contains(parsedMetadata, "variables")

	// Assert that the metadata file has 2 example variables: app_name and enable_https.
	assert.Equal("app_name", parsedMetadata["variables"].([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal("The name of the application.", parsedMetadata["variables"].([]interface{})[0].(map[string]interface{})["description"])
	assert.Equal("My App", parsedMetadata["variables"].([]interface{})[0].(map[string]interface{})["default"])
	assert.Equal("My App", parsedMetadata["variables"].([]interface{})[0].(map[string]interface{})["example"])

	assert.Equal("enable_https", parsedMetadata["variables"].([]interface{})[1].(map[string]interface{})["name"])
	assert.Equal("Whether to enable HTTPS or not.", parsedMetadata["variables"].([]interface{})[1].(map[string]interface{})["description"])
	assert.Equal(true, parsedMetadata["variables"].([]interface{})[1].(map[string]interface{})["example"])
}

// TestCreateCloneyProjectWithDefaultValues tests the creation of a new cloney project with default values,
// using the --non-interactive flag to force the use of default values.
func TestCreateCloneyProjectWithDefaultValues(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Simulate CLI arguments with flags and values.
	// Add the --non-interactive flag to force the use of default values.
	testStartCmd.SetArgs([]string{"--non-interactive"})

	// Execute the command.
	err := testStartCmd.Execute()

	// Assert that the command did not return an error.
	assert.Nil(err)

	// Assert that the command created the template repository directory.
	assert.DirExists(appConfig.DefaultCloneyProjectName)

	// Assert that the command created the metadata file.
	metadataFilePath := filepath.Join(appConfig.DefaultCloneyProjectName, appConfig.MetadataFileName)
	assert.FileExists(metadataFilePath)

	// Assert the metadata file contains the default values.
	metadataBytes, _ := os.ReadFile(metadataFilePath)
	var parsedMetadata map[string]interface{}
	err = yaml.Unmarshal(metadataBytes, &parsedMetadata)

	assert.Nil(err)
	assertMetadataMainFieldsAreOK(assert, TestMetadataProperties{
		Name:            appConfig.DefaultCloneyProjectName,
		Description:     appConfig.DefaultMetadataDescriptionValue,
		License:         appConfig.DefaultMetadataLicenseValue,
		TemplateVersion: appConfig.DefaultMetadataTemplateVersionValue,
		ManifestVersion: appConfig.MetadataManifestVersion,
	}, parsedMetadata)
	assertMetadataVariablesAreOK(assert, parsedMetadata)

	// Delete the created directory after the test.
	os.RemoveAll(appConfig.DefaultCloneyProjectName)

	// Reset the command flags.
	ResetStartCommandFlags(testStartCmd)
}

// TestCreateCloneyProjectWithFlags tests the creation of a new cloney project with flags and values.
func TestCreateCloneyProjectWithFlags(t *testing.T) {
	// Test the creation of a new cloney project using the project property flags, like --name, --description, etc.
	t.Run("TestCreateCloneyProjectWithProjectPropertiesFlags", func(t *testing.T) {
		// Create a new testing.T instance to use with assert functions.
		assert := assert.New(t)

		// Simulate CLI arguments with flags and values.
		testStartCmd.SetArgs([]string{
			"--name", "mock-name",
			"--description", "mock-description",
			"--license", "mock-license",
			"--non-interactive",
		})

		// Execute the command.
		err := testStartCmd.Execute()

		// Assert that the command did not return an error.
		assert.Nil(err)

		// Assert that the command created the template repository directory.
		assert.DirExists("mock-name")

		// Assert that the command created the metadata file.
		metadataFilePath := filepath.Join("mock-name", appConfig.MetadataFileName)
		assert.FileExists(metadataFilePath)

		// Assert the metadata file contains the default values.
		metadataBytes, _ := os.ReadFile(metadataFilePath)
		var metadataParsed map[string]interface{}
		err = yaml.Unmarshal(metadataBytes, &metadataParsed)

		assert.Nil(err)

		assertMetadataMainFieldsAreOK(assert, TestMetadataProperties{
			Name:            "mock-name",
			Description:     "mock-description",
			License:         "mock-license",
			TemplateVersion: appConfig.DefaultMetadataTemplateVersionValue,
			ManifestVersion: appConfig.MetadataManifestVersion,
		}, metadataParsed)
		assertMetadataVariablesAreOK(assert, metadataParsed)

		// Delete the created directory after the test.
		os.RemoveAll("mock-name")

		// Reset the command flags.
		ResetStartCommandFlags(testStartCmd)
	})

	// Test the creation of a new cloney project using the --output flag.
	t.Run("TestCreateCloneyProjectWithOutputFlag", func(t *testing.T) {
		// Create a new testing.T instance to use with assert functions.
		assert := assert.New(t)

		// Simulate CLI arguments with flags and values.
		testStartCmd.SetArgs([]string{
			"--output", "mock-output",
			"--non-interactive",
		})

		// Execute the command.
		err := testStartCmd.Execute()

		// Assert that the command did not return an error.
		assert.Nil(err)

		// Assert that the command created the template repository directory.
		assert.DirExists("mock-output")

		// Assert that the command created the metadata file.
		metadataFilePath := filepath.Join("mock-output", appConfig.MetadataFileName)
		assert.FileExists(metadataFilePath)

		// Assert the metadata file contains the default values.
		metadataBytes, _ := os.ReadFile(metadataFilePath)
		var metadataParsed map[string]interface{}
		err = yaml.Unmarshal(metadataBytes, &metadataParsed)

		assert.Nil(err)
		assertMetadataMainFieldsAreOK(assert, TestMetadataProperties{
			Name:            appConfig.DefaultCloneyProjectName,
			Description:     appConfig.DefaultMetadataDescriptionValue,
			License:         appConfig.DefaultMetadataLicenseValue,
			TemplateVersion: appConfig.DefaultMetadataTemplateVersionValue,
			ManifestVersion: appConfig.MetadataManifestVersion,
		}, metadataParsed)
		assertMetadataVariablesAreOK(assert, metadataParsed)

		// Delete the created directory after the test.
		os.RemoveAll("mock-output")

		// Reset the command flags.
		ResetStartCommandFlags(testStartCmd)
	})
}
