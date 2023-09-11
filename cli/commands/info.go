package commands

import (
	"fmt"

	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"

	"github.com/spf13/cobra"
)

// infoCmdRun is the function that runs when the info command is called.
func infoCmdRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Help()
		return
	}

	// Get arguments.
	repositoryURL := args[0]
	branch, _ := cmd.Flags().GetString("branch")

	// Create the repository struct.
	repository := &git.GitRepository{
		URL:    repositoryURL,
		Branch: branch,
	}

	// Validate arguments.
	err := repository.ValidateURL()
	if err != nil {
		fmt.Println("Invalid repository URL.")
		return
	}

	// Get the metadata file content.
	appConfig := config.GetAppConfig()
	metadataContent, err := repository.GetFileContent(appConfig.MetadataFileName)
	if err != nil {
		fmt.Println("Could not get repository metadata.", err)
		return
	}

	// Create the metadata struct.
	metadata, err := metadata.NewCloneyMetadata(metadataContent)
	if err != nil {
		fmt.Println("Could not parse repository metadata.", err)
		return
	}

	// Print metadata.
	metadata.PrettyPrint()

}

// infoCmd represents the info command.
// This command is used to print information about a Cloney template repository.
var infoCmd = &cobra.Command{
	Use:   "info [REPOSITORY_URL]",
	Short: "Prints information about a Cloney template repository",
	Run:   infoCmdRun,
}

// InitializeInfo initializes the info command.
func InitializeInfo(rootCmd *cobra.Command) {
	// Flags.
	infoCmd.Flags().StringP("branch", "b", "main", "Git branch")

	rootCmd.AddCommand(infoCmd)
}
