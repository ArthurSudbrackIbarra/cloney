package cli

import (
	"os"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

// rootCmdRun is the function that runs when the root command is called.
func rootCmdRun(cmd *cobra.Command, args []string) {
	// Display the root command's help information.
	cmd.Help()
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:          "cloney",
	Short:        "Cloney is a tool to clone template git repositories.",
	Run:          rootCmdRun,
	SilenceUsage: true,

	// Errors are printed by the commands.
	SilenceErrors: true,
}

// Initialize initializes the CLI.
func Initialize() {
	// Add subcommands.
	commands.InitializeVersion(rootCmd)
	commands.InitializeInfo(rootCmd)
	commands.InitializeClone(rootCmd)

	// CLI formatting, colors, bold, italic...
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		NoExtraNewlines: true,
	})

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		// Exit with an error code if there was an issue executing the command.
		os.Exit(1)
	}
}
