package commands

import (
	"fmt"

	git "github.com/ArthurSudbrackIbarra/cloney/git"

	cobra "github.com/spf13/cobra"
)

// cloneCmd represents the clone command.
// This command is used to clone a template repository.
var cloneCmd = &cobra.Command{
	Use:   "clone [REPOSITORY_URL]",
	Short: "Clones a template repository.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		repositoryURL := args[0]
		branch, _ := cmd.Flags().GetString("branch")
		path, _ := cmd.Flags().GetString("path")

		// Clone repository.
		fmt.Println("Cloning repository...")
		repository := &git.GitRepository{
			URL:    repositoryURL,
			Branch: branch,
		}
		err := git.CloneRepository(repository, path)
		if err != nil {
			fmt.Println(err)
			return
		}

	},
}

// InitializeClone initializes the clone command.
func InitializeClone(rootCmd *cobra.Command) {
	// Flags.
	cloneCmd.Flags().StringP("path", "p", "", "Path to clone the repository to.")
	cloneCmd.Flags().StringP("branch", "b", "main", "Branch to clone.")

	rootCmd.AddCommand(cloneCmd)
}
