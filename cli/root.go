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
	// Create subcommands.
	cloneCmd := commands.CreateCloneCommand()
	dryRunCmd := commands.CreateDryRunCommand()
	infoCmd := commands.CreateInfoCommand()
	startCmd := commands.CreateStartCommand()
	versionCmd := commands.CreateVersionCommand()

	// Add subcommands.
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(dryRunCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)

	// Stylings.
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.Bold,
		Commands:        cc.HiBlue + cc.Bold + cc.Italic,
		ExecName:        cc.HiYellow,
		Flags:           cc.HiRed + cc.Bold + cc.Italic,
		NoExtraNewlines: true,
	})

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		// Exit with an error code if there was an issue executing the command.
		os.Exit(1)
	}
}