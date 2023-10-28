package cli

import (
	"os"
	"strings"

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
	Short:        "Cloney is a tool to clone and create dynamic template git repositories.\nCopyright © - Arthur Sudbrack Ibarra - 2023.",
	Run:          rootCmdRun,
	SilenceUsage: true,

	// Errors are printed by the commands.
	SilenceErrors: true,
}

// Helper function to check if the error is related to "command not found."
func isUnknownCommandError(err error) bool {
	return strings.Contains(err.Error(), "unknown command")
}

// Helper function to check if the error is related to an "unknown flag."
func isUnknownFlagError(err error) bool {
	return strings.Contains(err.Error(), "unknown flag") || strings.Contains(err.Error(), "unknown shorthand flag")
}

// Initialize initializes the CLI.
func Initialize() {
	// Create subcommands.
	cloneCmd := commands.CreateCloneCommand()
	dryRunCmd := commands.CreateDryRunCommand()
	infoCmd := commands.CreateInfoCommand()
	startCmd := commands.CreateStartCommand()
	versionCmd := commands.CreateVersionCommand()
	validateCmd := commands.CreateValidateCommand()
	docsCmd := commands.CreateDocsCommand()

	// Add subcommands.
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(dryRunCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(docsCmd)

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
		// Print the error only if it is related to "command not found" or "unknown flag."
		if isUnknownCommandError(err) || isUnknownFlagError(err) {
			rootCmd.Println(err)
		}
		// Exit with an error code if there was an issue executing the command.
		os.Exit(1)
	}
}
