package cli

import (
	"fmt"
	"os"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands"

	"github.com/common-nighthawk/go-figure"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

// rootCmdPersistentPrePostRun is the function that runs before and after
// any command is run. This function is used to print new lines.
func rootCmdPersistentPrePostRun(cmd *cobra.Command, args []string) {
	fmt.Print("\n")
}

// rootCmdRun is the function that runs when the root command is called.
func rootCmdRun(cmd *cobra.Command, args []string) {
	cloneyASCIIArt := figure.NewFigure("cloney", "ogre", false)
	cloneyASCIIArt.Print()
	fmt.Print("\n")
	cmd.Help()
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:               "cloney",
	Short:             "Cloney is a tool to clone template git repositories.",
	PersistentPreRun:  rootCmdPersistentPrePostRun,
	PersistentPostRun: rootCmdPersistentPrePostRun,
	Run:               rootCmdRun,
	// Errors are printed by the commands.
	SilenceErrors: true,
}

// Initialize initializes the CLI.
func Initialize() {
	// Add subcommands.
	commands.InitializeVersion(rootCmd)
	commands.InitializeInfo(rootCmd)
	commands.InitializeClone(rootCmd)

	// Add colors to the CLI.
	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
