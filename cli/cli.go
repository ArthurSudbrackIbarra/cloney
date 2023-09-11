package cli

import (
	"fmt"
	"os"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands"

	"github.com/spf13/cobra"
)

// AddNewLine adds prints a new line to the terminal.
// This is used to add a new line before and after each command.
func AddNewLine(cmd *cobra.Command, args []string) {
	fmt.Print("\n")
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:              "cloney",
	Short:            "Cloney is a tool to clone template git repositories.",
	PersistentPreRun: AddNewLine,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPostRun: AddNewLine,
}

// Initialize initializes the CLI.
func Initialize() {
	// Add subcommands.
	commands.InitializeVersion(rootCmd)
	commands.InitializeInfo(rootCmd)
	commands.InitializeClone(rootCmd)

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
