package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/git"

	"github.com/spf13/cobra"
)

// cloneCmdRun is the function that runs when the clone command is called.
func cloneCmdRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Help()
		return
	}

	// Get arguments.
	repositoryURL := args[0]
	branch, _ := cmd.Flags().GetString("branch")
	path, _ := cmd.Flags().GetString("path")

	// Validate arguments.
	err := git.ValidateRepositoryURL(repositoryURL)
	if err != nil {
		fmt.Println("Invalid repository URL.")
		return
	}

	// Get user current directory.
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get user's current directory.")
		return
	}
	path = filepath.Join(currentDir, path)

	// Clone repository.
	repository := &git.GitRepository{
		URL:    repositoryURL,
		Branch: branch,
	}
	repositoryName := git.GetRepositoryName(repositoryURL)
	fmt.Printf("Cloning '%s' into '%s' ...\n", repositoryName, path)
	err = git.CloneRepository(repository, path)
	if err != nil {
		fmt.Println("Could not clone repository.", err)
		return
	}

}

// cloneCmd represents the clone command.
// This command is used to clone a template repository.
var cloneCmd = &cobra.Command{
	Use:   "clone [REPOSITORY_URL]",
	Short: "Clones a template repository.",
	Run:   cloneCmdRun,
}

// InitializeClone initializes the clone command.
func InitializeClone(rootCmd *cobra.Command) {
	// Flags.
	cloneCmd.Flags().StringP("path", "p", "", "Path to clone the repository to.")
	cloneCmd.Flags().StringP("branch", "b", "main", "Branch to clone.")

	rootCmd.AddCommand(cloneCmd)
}
