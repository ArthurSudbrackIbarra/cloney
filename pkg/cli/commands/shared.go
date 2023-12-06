package commands

import (
	"github.com/ArthurSudbrackIbarra/cloney/pkg/config"
	"github.com/ArthurSudbrackIbarra/cloney/pkg/terminal"

	"github.com/spf13/cobra"
)

// appConfig stores the application configuration.
// This configuration is retrieved once and shared across all commands.
var appConfig = config.GetAppConfig()

// persistentPreRun is executed before every command.
// It sets the current command, allowing messages to be printed to the command's output using cmd.Print().
// This is essential for capturing the output of the command to be used in the tests.
func persistentPreRun(cmd *cobra.Command, args []string) {
	terminal.SetCmd(cmd)
}
