package cli

import (
	"fmt"
	"os"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// rootCmdPersistentPrePostRun is the function that runs before and after
// any command is run. This function is used to add a new lines.
func rootCmdPersistentPrePostRun(cmd *cobra.Command, args []string) {
	fmt.Print("\n")
}

// rootCmdRun is the function that runs when the root command is called.
func rootCmdRun(cmd *cobra.Command, args []string) error {
	cloneyArt := figure.NewFigure("cloney", "doom", true)
	cloneyArt.Print()
	fmt.Print("\n")
	cmd.Help()
	return nil
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:               "cloney",
	Short:             "Cloney is a tool to clone template git repositories.",
	PersistentPreRun:  rootCmdPersistentPrePostRun,
	PersistentPostRun: rootCmdPersistentPrePostRun,
	RunE:              rootCmdRun,
}

// Initialize initializes the CLI.
func Initialize() {
	// Add subcommands.
	commands.InitializeVersion(rootCmd)
	commands.InitializeInfo(rootCmd)
	commands.InitializeClone(rootCmd)

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
