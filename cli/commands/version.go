package commands

import (
	"fmt"

	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/spf13/cobra"
)

// versionCmdRun is the function that runs when the version command is called.
func versionCmdRun(cmd *cobra.Command, args []string) {
	appConfig := config.GetAppConfig()
	fmt.Printf("Cloney version %s\n", appConfig.Version)
}

// versionCmd represents the version command.
// This command is used to print the version of the application.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of Cloney",
	Long:  "\nPrints the current version of Cloney.",
	Run:   versionCmdRun,
}

// InitializeVersion initializes the version command.
func InitializeVersion(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}
