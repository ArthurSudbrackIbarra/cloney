package commands

import (
	"fmt"

	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"

	"github.com/spf13/cobra"
)

// infoCmdRun is the function that runs when the info command is called.
func infoCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		// Display command help if no repository URL is provided.
		cmd.Help()
		return nil
	}

	// Get command-line arguments.
	repositoryURL := args[0]
	branch, _ := cmd.Flags().GetString("branch")

	// Create the Git repository instance.
	repository, err := git.NewGitRepository(repositoryURL, branch)
	if err != nil {
		// Handle errors related to the repository.
		fmt.Println("Error referencing repository:", err)
		return err
	}

	// Get the metadata file content.
	appConfig := config.GetAppConfig()
	metadataContent, err := repository.GetFileContent(appConfig.MetadataFileName)
	if err != nil {
		// Handle errors related to reading the metadata file.
		fmt.Println(
			fmt.Sprintf("Error reading the repository '%s' metadata file:", appConfig.MetadataFileName), err,
		)
		return err
	}

	// Create the metadata struct from raw YAML data.
	metadata, err := metadata.NewCloneyMetadataFromRawYAML(metadataContent)
	if err != nil {
		// Handle errors related to parsing repository metadata.
		fmt.Println("Could not parse repository metadata file:", err)
		return err
	}

	// Print metadata.
	metadata.Show()
	return nil
}

// infoCmd represents the info command.
// This command is used to print information about a Cloney template repository.
var infoCmd = &cobra.Command{
	Use:     "info [repository_url]",
	Short:   "Prints information about a Cloney template repository",
	Long:    "\nPrints information about a Cloney template repository.",
	Example: "  cloney info https://github.com/ArthurSudbrackIbarra/cloney.git",
	RunE:    infoCmdRun,
}

// InitializeInfo initializes the info command.
func InitializeInfo(rootCmd *cobra.Command) {
	// Define command-line flags.
	infoCmd.Flags().StringP("branch", "b", "main", "Git branch")

	// Add the info command to the root command.
	rootCmd.AddCommand(infoCmd)
}
