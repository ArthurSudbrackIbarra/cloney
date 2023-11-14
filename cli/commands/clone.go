package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/terminal"

	"github.com/spf13/cobra"
)

// cloneCmdRun is the function that runs when the 'clone' command is called.
func cloneCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		terminal.ErrorMessage("You must provide a repository URL\n", nil)

		// Display command help if no repository URL is provided.
		cmd.Help()
		return nil
	}

	// Get command-line arguments.
	repositoryURL := args[0]
	branch, _ := cmd.Flags().GetString("branch")
	output, _ := cmd.Flags().GetString("output")
	tag, _ := cmd.Flags().GetString("tag")
	variables, _ := cmd.Flags().GetString("variables")
	token, _ := cmd.Flags().GetString("token")

	// Variable to store errors.
	var err error

	// Get the current working directory.
	currentDir, err := steps.GetCurrentWorkingDirectory()
	if err != nil {
		return err
	}

	// Get the template variables provided by the user.
	variablesMap, err := steps.GetUserVariablesMap(currentDir, variables)
	if err != nil {
		return err
	}

	// Create and validate the Git repository.
	repository, err := steps.CreateAndValidateRepository(repositoryURL, branch, tag)
	if err != nil {
		return err
	}

	// If a token is provided, authenticate with it.
	steps.AuthenticateToRepository(repository, token)

	// Calculate the clone path.
	clonePath, _ := steps.CalculatePath(output, repository.GetName())

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

	// Delete the .git directory.
	gitDirPath := filepath.Join(clonePath, ".git")
	os.RemoveAll(gitDirPath)

	// Parse the metadata file.
	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent, appConfig.SupportedManifestVersions)
	if err != nil {
		// If it was not possible to parse the metadata file, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	// Validate if the user variables match the template variables.
	// Also, fill default values of the variables if they are not defined.
	err = steps.MatchUserVariables(cloneyMetadata, variablesMap)
	if err != nil {
		// If the user variables do not match the template variables, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	// Define options for ignoring specific files and directories when filling template variables.
	var ignorePaths []string
	ignorePaths = append(ignorePaths, appConfig.KnownIgnorePaths...)
	ignorePaths = append(ignorePaths, cloneyMetadata.Configuration.IgnorePaths...)

	// Set the 'outputInTerminal' parameter to 'false' because we intend to actually fill the template variables.
	err = steps.FillDirectory(clonePath, ignorePaths, false, variablesMap)
	if err != nil {
		// If it was not possible to fill the template variables, delete the cloned repository.
		os.RemoveAll(clonePath)
		return err
	}

	// Delete the paths specified in the 'ignore_paths' field of the metadata file.
	steps.DeleteIgnoredPaths(clonePath, ignorePaths)

	terminal.Message("\nDone!")

	return nil
}

// CreateCloneCommand creates the 'clone' command and its respective flags.
func CreateCloneCommand() *cobra.Command {
	// cloneCmd represents the 'clone' command.
	// This command is used to clone a template repository.
	cloneCmd := &cobra.Command{
		Use:   "clone [repository_url]",
		Short: "Clone a template repository",
		Long: fmt.Sprintf(`Clone a template repository.

The 'cloney clone' command will search for a file named '%s' in your current directory by default.
You can specify a different file or pass the variables inline as YAML using the '--variables' flag.`, appConfig.DefaultUserVariablesFileName),
		Example: strings.Join([]string{
			"  clone https://github.com/username/repository.git",
			"  clone https://github.com/username/repository.git -v variables.yaml",
			"  clone https://github.com/username/repository.git -v '{ var1: value, var2: value }'",
		}, "\n"),
		Aliases:          []string{"cl"},
		PersistentPreRun: persistentPreRun,
		RunE:             cloneCmdRun,
	}

	// Define command-line flags for the 'clone' command.
	cloneCmd.Flags().StringP("output", "o", "", "Path to clone the repository to")
	cloneCmd.Flags().StringP("branch", "b", "main", "Git branch")
	cloneCmd.Flags().StringP("tag", "t", "", "Git tag")
	cloneCmd.Flags().StringP("variables", "v", appConfig.DefaultUserVariablesFileName, "Path to a template variables file or raw YAML")
	cloneCmd.Flags().StringP("token", "k", "", "Git token, if referencing a private Git repository (not recommended)")

	return cloneCmd
}
