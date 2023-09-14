package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/templates"

	"github.com/spf13/cobra"
)

// cloneCmdRun is the function that runs when the clone command is called.
func cloneCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		// Display command help if no repository URL is provided.
		cmd.Help()
		return nil
	}

	// Get command-line arguments.
	repositoryURL := args[0]
	branch, _ := cmd.Flags().GetString("branch")
	output, _ := cmd.Flags().GetString("output")
	tag, _ := cmd.Flags().GetString("tag")
	variablesFilePath, _ := cmd.Flags().GetString("variables-file")
	variablesJSON, _ := cmd.Flags().GetString("variables")

	// Variable to store errors.
	var err error

	// Get the current working directory.
	currentDir, err := steps.GetCurrentWorkingDirectory()
	if err != nil {
		return err
	}

	// Get the template variables provided by the user.
	variablesMap, err := steps.GetUserVariablesMap(currentDir, variablesJSON, variablesFilePath)
	if err != nil {
		return err
	}

	// Create and validate the git repository.
	repository, err := steps.CreateAndValidateRepository(repositoryURL, branch, tag, output)

	// If a token is provided, authenticate with it.
	appConfig := config.GetAppConfig()
	steps.AuthenticateToRepository(repository, appConfig.GitToken)

	// Calculate the clone path.
	clonePath := steps.CalculateClonePath(repository, currentDir, output)

	// Clone the repository.
	err = steps.CloneRepository(repository, clonePath)
	if err != nil {
		return err
	}

	// Read the repository metadata file.
	metadataFilePath := filepath.Join(clonePath, appConfig.MetadataFileName)
	metadataContent, err := steps.ReadRepositoryMetadata(metadataFilePath)
	if err != nil {
		// If it was not possible to read the metadata file, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	// Delete the repository metadata file.
	os.Remove(metadataFilePath)

	// Parse the metadata file.
	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent)
	if err != nil {
		// If it was not possible to parse the metadata file, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	// Validate if the user variables match the template variables.
	// Also fill default values of the variables if they are not defined.
	err = steps.MatchUserVariables(cloneyMetadata, variablesMap)
	if err != nil {
		// If the user variables do not match the template variables, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	// Fill the template variables in the cloned directory.
	options := templates.TemplateFillOptions{
		SourceDirectoryPath: clonePath,
	}
	err = steps.FillTemplateVariables(options, variablesMap)
	if err != nil {
		// If it was not possible to fill the template variables, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	fmt.Println("\nDone!")
	return nil
}

// cloneCmd represents the clone command.
// This command is used to clone a template repository.
var cloneCmd = &cobra.Command{
	Use:     "clone [repository_url]",
	Short:   "Clones a template repository.",
	Long:    "\nClones a template repository.",
	Example: "  cloney clone https://github.com/ArthurSudbrackIbarra/cloney.git",
	Aliases: []string{"cl"},
	RunE:    cloneCmdRun,
}

// InitializeClone initializes the clone command.
func InitializeClone(rootCmd *cobra.Command) {
	appConfig := config.GetAppConfig()

	// Define command-line flags.
	cloneCmd.Flags().StringP("output", "o", "", "Path to clone the repository to")
	cloneCmd.Flags().StringP("branch", "b", "main", "Git branch")
	cloneCmd.Flags().StringP("tag", "t", "", "Git tag")
	cloneCmd.Flags().StringP("variables-file", "f", appConfig.DefaultUserVariablesFileName, "Path to a template variables YAML or JSON file")
	cloneCmd.Flags().StringP("variables", "v", "", "Inline template variables as JSON")

	// Add the clone command to the root command.
	rootCmd.AddCommand(cloneCmd)
}
