package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/git"

	"github.com/spf13/cobra"
)

// cloneCmdRun is the function that runs when the clone command is called.
func cloneCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Help()
		return nil
	}

	// Get arguments.
	repositoryURL := args[0]
	branch, _ := cmd.Flags().GetString("branch")
	path, _ := cmd.Flags().GetString("path")

	// Create the repository struct.
	repository, err := git.NewGitRepository(repositoryURL, branch)
	if err != nil {
		fmt.Println("Error referencing repository.", err)
		return err
	}

	// Get user current directory.
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get user's current directory.")
		return err
	}
	path = filepath.Join(currentDir, path)

	// Clone repository.
	repositoryName := repository.GetName()
	fmt.Printf("Cloning '%s' into '%s' ...\n", repositoryName, path)
	err = repository.Clone(path)
	if err != nil {
		fmt.Println("Could not clone repository.", err)
		return err
	}

	return nil
}

// cloneCmd represents the clone command.
// This command is used to clone a template repository.
var cloneCmd = &cobra.Command{
	Use:     "clone [repository_url]",
	Short:   "Clones a template repository.",
	Example: "  cloney clone https://github.com/ArthurSudbrackIbarra/cloney.git",
	RunE:    cloneCmdRun,
}

// InitializeClone initializes the clone command.
func InitializeClone(rootCmd *cobra.Command) {
	// Flags.
	cloneCmd.Flags().StringP("path", "p", "", "Path to clone the repository to")
	cloneCmd.Flags().StringP("branch", "b", "main", "Git branch")

	rootCmd.AddCommand(cloneCmd)
}
