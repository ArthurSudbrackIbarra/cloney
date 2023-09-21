package commands

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
// 	"github.com/spf13/cobra"

// 	"github.com/stretchr/testify/assert"
// )

// // TestStartCommandNonInteractive tests the startCmdRun function with the --non-interactive flag.
// // This test case validates that if the user uses the --non-interactive flag, the Cloney project
// // must be created in the default path and the metadata file must be filled with the default values.
// func TestStartCommandNonInteractive(t *testing.T) {
// 	// Call the startCmdRun function with the test parameters.
// 	startCmd.Flags().Set("non-interactive", "true")
// 	err := startCmdRun(startCmd, []string{})

// 	// No error should be returned.
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	// Check if the Cloney project was created in the default path.
// 	projectPath, err := steps.CalculatePath(
// 		appConfig.DefaultCloneyProjectName,
// 		"",
// 	)
// 	if err != nil {
// 		t.Errorf("Project should have been created, got %v", err)
// 	}

// 	// Read the repository metadata file.
// 	metadataFilePath := filepath.Join(projectPath, appConfig.MetadataFileName)
// 	metadataContent, err := steps.ReadRepositoryMetadata(metadataFilePath)
// 	if err != nil {
// 		t.Errorf("Metadata file should have been created, got %v", err)
// 	}

// 	// Check if the metadata file contains the expected content.
// 	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent, appConfig.SupportedManifestVersions)
// 	if err != nil {
// 		t.Errorf("Metadata file should be able to be parsed, got %v", err)
// 	}
// 	assert.Equal(t, appConfig.DefaultCloneyProjectName, cloneyMetadata.Name)
// 	assert.Equal(t, appConfig.DefaultMetadataDescriptionValue, cloneyMetadata.Description)
// 	assert.Equal(t, appConfig.DefaultMetadataLicenseValue, cloneyMetadata.License)
// 	assert.Equal(t, appConfig.MetadataManifestVersion, cloneyMetadata.ManifestVersion)
// 	assert.Equal(t, appConfig.DefaultMetadataTemplateVersionValue, cloneyMetadata.TemplateVersion)

// 	// Delete the Cloney project after the test.
// 	os.RemoveAll(projectPath)
// }

// // TestStartCommandWhenProjectFlagsAreUsed tests the startCmdRun function when the user use flags to set the project properties,
// // like --name, --description, --license...
// // Check if the metadata file contains the expected content, which is the same as the flags.
// func TestStartCommandWhenProjectFlagsAreUsed(t *testing.T) {
// 	// Create a new Cobra command for testing.
// 	cmd := &cobra.Command{
// 		Use: "test",
// 	}

// 	// Call the startCmdRun function with the test parameters.
// 	cmd.Flags().Set("name", "name")
// 	cmd.Flags().Set("description", "description")
// 	cmd.Flags().Set("license", "license")
// 	cmd.Flags().Set("authors", "authors_1, authors_2")
// 	cmd.Flags().Set("non-interactive", "true")
// 	err := startCmdRun(cmd, []string{})

// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	// Check if the Cloney project was created in the default path.
// 	projectPath, err := steps.CalculatePath(
// 		appConfig.DefaultCloneyProjectName,
// 		"",
// 	)
// 	if err != nil {
// 		t.Errorf("Project should have been created, got %v", err)
// 	}

// 	// Read the repository metadata file.
// 	metadataFilePath := filepath.Join(projectPath, appConfig.MetadataFileName)
// 	metadataContent, err := steps.ReadRepositoryMetadata(metadataFilePath)
// 	if err != nil {
// 		t.Errorf("Metadata file should have been created, got %v", err)
// 	}

// 	fmt.Println(metadataContent)

// 	// Check if the metadata file contains the expected content.
// 	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent, appConfig.SupportedManifestVersions)
// 	if err != nil {
// 		t.Errorf("Metadata file should be able to be parsed, got %v", err)
// 	}
// 	assert.Equal(t, "name", cloneyMetadata.Name)
// 	assert.Equal(t, "description", cloneyMetadata.Description)
// 	assert.Equal(t, "license", cloneyMetadata.License)
// 	assert.Equal(t, appConfig.MetadataManifestVersion, cloneyMetadata.ManifestVersion)
// 	for i, author := range cloneyMetadata.Authors {
// 		assert.Equal(t, fmt.Sprintf("authors_%d", i+1), author)
// 	}

// 	// Delete the Cloney project after the test.
// 	os.RemoveAll(projectPath)
// }
