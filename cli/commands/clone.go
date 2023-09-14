package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"
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
	path, _ := cmd.Flags().GetString("path")
	tag, _ := cmd.Flags().GetString("tag")
	variablesFilePath, _ := cmd.Flags().GetString("variables-file")
	variablesJSON, _ := cmd.Flags().GetString("variables")

	// Create the Git repository instance.
	repository := &git.GitRepository{
		URL:    repositoryURL,
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
		repository.AuthenticateWithToken(appConfig.GitToken)
	}

	// Get the name of the repository.
	repositoryName := repository.GetName()

	// Get the current working directory.
	currentDir, err := os.Getwd()
	if err != nil {
		// Handle errors related to getting the current user's directory.
		fmt.Println("Could not get user's current directory:", err)
		return err
	}

	// Calculate the clone path.
	clonePath := filepath.Join(currentDir, path, repositoryName)

	// Clone the repository.
	err = repository.Clone(clonePath)
	if err != nil {
		// Handle errors related to cloning the repository.
		fmt.Println("Could not clone repository:", err)
		return err
	}

	// Read the repository metadata file.
	metadataBytes, err := os.ReadFile(
		filepath.Join(clonePath, appConfig.MetadataFileName),
	)
	if err != nil {
		// Handle errors related to reading the metadata file.
		fmt.Printf("Could not read repository '%s' metadata file: %s\n", appConfig.MetadataFileName, err)
		return err
	}

	// Create the metadata struct from raw YAML data.
	cloneyMetadata, err := metadata.NewCloneyMetadataFromRawYAML(string(metadataBytes))
	if err != nil {
		// Handle errors related to parsing repository metadata.
		fmt.Println("Could not parse repository metadata:", err)
		return err
	}

	// Get the template variables provided by the user.
	var variablesMap map[string]interface{}
	if variablesJSON != "" {
		variablesMap, err = metadata.NewCloneyUserVariablesFromRawJSON(variablesJSON)
		if err != nil {
			// Handle errors related to parsing user template variables.
			fmt.Println("Could not parse template variables:", err)
			return err
		}
	} else {
		variablesMap, err = metadata.NewCloneyUserVariablesFromFile(variablesFilePath)
		if err != nil {
			// Handle errors related to reading user template variables.
			fmt.Println("Could not read template variables:", err)
			return err
		}
	}

	// Validate if the user variables match the template variables.
	// Also fill default values of the variables if they are not defined.
	variablesMap, err = cloneyMetadata.MatchUserVariables(variablesMap)
	if err != nil {
		// Handle errors related to validating template variables.
		fmt.Println("Error validating template variables:", err)
		// Delete the cloned directory.
		os.RemoveAll(clonePath)
		return err
	}

	// Fill the template variables in the cloned directory.
	filler := templates.NewTemplateFiller(variablesMap)
	err = filler.FillDirectory(clonePath)
	if err != nil {
		// Handle errors related to filling template variables.
		fmt.Println("Error filling template variables:", err)
		// Delete the cloned directory.
		os.RemoveAll(clonePath)
		return err
	}

	return nil
}

// cloneCmd represents the clone command.
// This command is used to clone a template repository.
var cloneCmd = &cobra.Command{
	Use:     "clone [repository_url]",
	Short:   "Clones a template repository.",
	Long:    "\nClones a template repository.",
	Example: "  cloney clone https://github.com/ArthurSudbrackIbarra/cloney.git",
	RunE:    cloneCmdRun,
}

// InitializeClone initializes the clone command.
func InitializeClone(rootCmd *cobra.Command) {
	appConfig := config.GetAppConfig()

	// Define command-line flags.
	cloneCmd.Flags().StringP("path", "p", "", "Path to clone the repository to")
	cloneCmd.Flags().StringP("branch", "b", "main", "Git branch")
	cloneCmd.Flags().StringP("tag", "t", "", "Git tag")
	cloneCmd.Flags().StringP("variables-file", "f", appConfig.DefaultUserVariablesFileName, "Path to a template variables YAML or JSON file")
	cloneCmd.Flags().StringP("variables", "v", "", "Inline template variables as JSON")

	// Add the clone command to the root command.
	rootCmd.AddCommand(cloneCmd)
}
