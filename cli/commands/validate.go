package commands

import (
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"

	"github.com/spf13/cobra"
)

// validateCmd is the function that runs when the 'validate' command is called.
func validateCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	path, _ := cmd.Flags().GetString("path")

	// Variable to store errors.
	var err error

	// Calculate the template directory path.
	sourcePath, err := steps.CalculatePath(path, "")
	if err != nil {
		return err
	}

	// Read the repository metadata file.
	metadataFilePath := filepath.Join(sourcePath, appConfig.MetadataFileName)
	metadataContent, err := steps.ReadRepositoryMetadata(metadataFilePath)
	if err != nil {
		return err
	}

	// Parse the metadata file.
	_, err = steps.ParseRepositoryMetadata(metadataContent, appConfig.SupportedManifestVersions)
	if err != nil {
		return err
	}

	// If the metadata file was parsed successfully, then the template is valid.
	cmd.Println("\nYour Cloney template is valid!")

	return nil
}

// ResetValidateCommandFlags resets the flags of the 'validate' command.
func ResetValidateCommandFlags(cmd *cobra.Command) {
	cmd.Flags().Set("path", "")
}

// CreateValidateCommand creates the 'validate' command.
func CreateValidateCommand() *cobra.Command {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate your Cloney template repository",
		Long: `Validate your Cloney template repository.

The 'cloney validate' command validates if your Cloney template repository is valid.
It checks if the repository has a metadata file, and if it has the required fields in it.
`,
		PersistentPreRun: persistentPreRun,
		RunE:             validateCmdRun,
	}

	// Define command-line flags for the 'validate' command.
	validateCmd.Flags().StringP("path", "p", "", "Path to your local template repository")

	return validateCmd
}
