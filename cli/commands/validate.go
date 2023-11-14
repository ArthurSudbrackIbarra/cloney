package commands

import (
	"path/filepath"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/terminal"

	"github.com/spf13/cobra"
)

// validateCmd is the function that runs when the 'validate' command is called.
func validateCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	var repositorySource string
	if len(args) >= 1 {
		repositorySource = args[0]
	}

	// Variable to store errors.
	var err error

	// Calculate the template directory path.
	sourcePath, err := steps.CalculatePath(repositorySource, "")
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
	terminal.Message("\nYour Cloney template is valid!")

	return nil
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
		Example: strings.Join([]string{
			"  validate",
			"  validate ./path/to/my/template",
		}, "\n"),
		PersistentPreRun: persistentPreRun,
		RunE:             validateCmdRun,
	}

	return validateCmd
}
