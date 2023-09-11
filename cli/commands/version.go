package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command.
// This command is used to print the version of the application.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of Cloney",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloney v0.1.0")
	},
}

// InitializeVersion initializes the version command.
func InitializeVersion(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}
