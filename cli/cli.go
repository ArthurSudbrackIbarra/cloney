package cli

import (
	"fmt"
	"os"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "cloney",
	Short: "Cloney is a tool to clone template git repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Initialize initializes the CLI.
func Initialize() {
	// Add subcommands.
	commands.InitializeVersion(rootCmd)
	commands.InitializeClone(rootCmd)

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
