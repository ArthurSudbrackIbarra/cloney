package commands

import (
	"fmt"

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
	metadataContent, err := repository.GetFileContent(".cloney.yaml")
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
	fmt.Println("Name:", metadata.Name)
	fmt.Println("Description:", metadata.Description)
	fmt.Println("Version:", metadata.Version)
	fmt.Println("Authors:", metadata.Authors)
	fmt.Println("License:", metadata.License)
	fmt.Println("Variables:")
	for _, variable := range metadata.Variables {
		fmt.Println("\tName:", variable.Name)
		fmt.Println("\tDescription:", variable.Description)
		fmt.Println("\tDefault:", variable.Default)
		fmt.Println("\tType:", variable.Type)
		fmt.Println()
	}

}

// infoCmd represents the info command.
// This command is used to print information about a Cloney template repository.
var infoCmd = &cobra.Command{
	Use:   "info [REPOSITORY_URL]",
	Short: "Prints information about a Cloney template repository.",
	Run:   infoCmdRun,
}

// InitializeInfo initializes the info command.
func InitializeInfo(rootCmd *cobra.Command) {
	// Flags.
	infoCmd.Flags().StringP("branch", "b", "main", "Git branch.")

	rootCmd.AddCommand(infoCmd)
}
