package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/terminal"

	"github.com/spf13/cobra"
)

// infoCmdRun is the function that runs when the 'info' command is called.
func infoCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	var repositorySource string
	if len(args) >= 1 {
		repositorySource = args[0]
	}
	branch, _ := cmd.Flags().GetString("branch")
	tag, _ := cmd.Flags().GetString("tag")
	token, _ := cmd.Flags().GetString("token")

	// Variable to store errors.
	var err error

	// Suppress prints for this command.
	steps.SetSuppressPrints(true)

	// Variable to store the metadata file content.
	var metadataContent string

	// If the argument is a git repository URL, use it.
	if git.MatchesGitRepositoryURL(repositorySource) {
		// Create and validate the git repository.
		repository, err := steps.CreateAndValidateRepository(repositorySource, branch, tag)
		if err != nil {
			return err
		}

		// If a token is provided, authenticate with it.
		if token == "" {
			token = appConfig.GitToken
		}
		steps.AuthenticateToRepository(repository, token)

		// Get the metadata file content.
		metadataContent, err = repository.GetFileContent(appConfig.MetadataFileName)
		if err != nil {
			terminal.ErrorMessage(
				fmt.Sprintf("Error reading the repository '%s' metadata file:", appConfig.MetadataFileName), err,
			)
			return err
		}
	} else {
		// If the argument is not a git repository URL, assume it is a local path.

		// Calculate the directory path.
		repositorySource, _ := steps.CalculatePath(repositorySource, "")

		// Get the metadata file content.
		metadataFilePath := filepath.Join(repositorySource, config.GetAppConfig().MetadataFileName)
		metadataContent, err = steps.ReadRepositoryMetadata(metadataFilePath)
		if err != nil {
			return err
		}
	}

	// Create the metadata struct from raw YAML data.
	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent, appConfig.SupportedManifestVersions)
	if err != nil {
		return err
	}

	// Print metadata.
	cmd.Println(cloneyMetadata)

	return nil
}

// ResetInfoCommandFlags resets the flags of the 'info' command.
func ResetInfoCommandFlags(cmd *cobra.Command) {
	cmd.Flags().Set("branch", "main")
	cmd.Flags().Set("tag", "")
	cmd.Flags().Set("token", "")
}

// CreateInfoCommand creates the 'info' command and its respective flags.
func CreateInfoCommand() *cobra.Command {
	// infoCmd represents the info command.
	// This command is used to print information about a Cloney template repository.
	infoCmd := &cobra.Command{
		Use:   "info [local_path OR repository_url]",
		Short: "Get information about a Cloney template repository",
		Long: `Get information about a Cloney template repository.

The 'cloney info' command will give you information about a Cloney template repository, such as its name, description, and variables.

It can get information from a local template repository, or from a remote template repository.
By default, it will get information from the current directory, assuming it is a template repository.
`,
		Example: strings.Join([]string{
			"  info",
			"  info ./path/to/my/template",
			"  info https://github.com/username/repository.git",
		}, "\n"),
		Aliases:          []string{"more"},
		PersistentPreRun: persistentPreRun,
		RunE:             infoCmdRun,
	}

	// Define command-line flags for the 'info' command.
	infoCmd.Flags().StringP("branch", "b", "main", "Git branch, if referencing a git repository")
	infoCmd.Flags().StringP("tag", "t", "", "Git tag, if referencing a git repository")
	infoCmd.Flags().StringP("token", "k", "", "Git token, if referencing a private git repository")

	return infoCmd
}
