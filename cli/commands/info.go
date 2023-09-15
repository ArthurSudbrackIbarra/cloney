package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"

	"github.com/spf13/cobra"
)

// infoCmdRun is the function that runs when the info command is called.
func infoCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	var repositorySource string
	if len(args) >= 1 {
		repositorySource = args[0]
	}
	branch, _ := cmd.Flags().GetString("branch")
	tag, _ := cmd.Flags().GetString("tag")

	// Variable to store the metadata file content.
	var metadataContent string

	// If the argument is a git repository URL, use it.
	if git.MatchesGitRepositoryURL(repositorySource) {
		// Create the Git repository instance.
		repository := &git.GitRepository{
			URL:    repositorySource,
			Branch: branch,
			Tag:    tag,
		}

		// Validate the repository.
		err := repository.Validate()
		if err != nil {
			// Handle errors related to the repository.
			fmt.Println("Error validating repository:", err)
			return err
		}

		// If a token is provided, authenticate with it.
		appConfig := config.GetAppConfig()
		if appConfig.GitToken != "" {
			fmt.Println("Authenticating with token...")
			fmt.Println(appConfig.GitToken)
			repository.AuthenticateWithToken(appConfig.GitToken)
		}

		// Get the metadata file content.
		metadataContent, err = repository.GetFileContent(appConfig.MetadataFileName)
		if err != nil {
			// Handle errors related to reading the metadata file.
			fmt.Println(
				fmt.Sprintf("Error reading the repository '%s' metadata file:", appConfig.MetadataFileName), err,
			)
			return err
		}
	} else {
		// If the argument is not a git repository URL, assume it is a local path.

		// Get the current working directory.
		currentDir, err := steps.GetCurrentWorkingDirectory()
		if err != nil {
			return err
		}

		// Get the metadata file content.
		steps.SetSuppressPrints(true)
		metadataFilePath := filepath.Join(currentDir, repositorySource, config.GetAppConfig().MetadataFileName)
		metadataContent, err = steps.ReadRepositoryMetadata(metadataFilePath)
		if err != nil {
			return err
		}
		steps.SetSuppressPrints(false)
	}

	// Create the metadata struct from raw YAML data.
	cloneyMetadata, err := metadata.NewCloneyMetadataFromRawYAML(metadataContent)
	if err != nil {
		// Handle errors related to parsing repository metadata.
		fmt.Println("Could not parse repository metadata file:", err)
		return err
	}

	// Print metadata.
	cloneyMetadata.Show()
	return nil
}

// infoCmd represents the info command.
// This command is used to print information about a Cloney template repository.
var infoCmd = &cobra.Command{
	Use:   "info [local_path OR repository_url]",
	Short: "Prints information about a Cloney template repository",
	Long:  "\nPrints information about a Cloney template repository.",
	Example: strings.Join([]string{
		"  cloney info https://github.com/ArthurSudbrackIbarra/cloney.git -- Info about a remote template repository",
		"  cloney info ./my-template -- Info about a local template repository in the given path",
		"  cloney info -- Info about the current directory as a template repository",
	}, "\n"),
	Aliases: []string{"more"},
	RunE:    infoCmdRun,
}

// InitializeInfo initializes the info command.
func InitializeInfo(rootCmd *cobra.Command) {
	// Define command-line flags.
	infoCmd.Flags().StringP("branch", "b", "main", "Git branch, if referencing a git repository")
	infoCmd.Flags().StringP("tag", "t", "", "Git tag, if referencing a git repository")

	// Add the info command to the root command.
	rootCmd.AddCommand(infoCmd)
}
