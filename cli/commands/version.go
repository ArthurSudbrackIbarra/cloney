package commands

import (
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmdRun is the function that runs when the 'version' command is called.
func versionCmdRun(cmd *cobra.Command, args []string) {
	// Print the current version of Cloney, the operating system and the architecture.
	cmd.Printf("Cloney version %s %s %s\n", appConfig.AppVersion, runtime.GOOS, runtime.GOARCH)
}

// CreateVersionCommand creates the 'version' command.
func CreateVersionCommand() *cobra.Command {
	// versionCmd represents the version command.
	// This command is used to print the version of the application.
	versionCmd := &cobra.Command{
		Use:              "version",
		Short:            "Prints the current version of Cloney",
		PersistentPreRun: persistentPreRun,
		Run:              versionCmdRun,
	}

	return versionCmd
}
