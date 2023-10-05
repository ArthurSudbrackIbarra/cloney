package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/ArthurSudbrackIbarra/cloney/terminal"
	"github.com/spf13/cobra"
)

// docsCmd is the function that runs when the 'docs' command is called.
func docsCmdRun(cmd *cobra.Command, args []string) error {
	// Variable to store errors.
	var err error

	// Open the Cloney documentation in the default browser.
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", appConfig.CloneyDocumentationURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", appConfig.CloneyDocumentationURL).Start()
	case "darwin":
		err = exec.Command("open", appConfig.CloneyDocumentationURL).Start()
	default:
		terminal.ErrorMessage("Compatibility error", fmt.Errorf("unsupported operating system"))
	}
	if err != nil {
		terminal.ErrorMessage("Error opening Cloney documentation in default browser", err)
	}

	return nil
}

// CreateDocsCommand creates the 'docs' command.
func CreateDocsCommand() *cobra.Command {
	docsCmd := &cobra.Command{
		Use:              "docs",
		Short:            "Open the Cloney documentation in your browser",
		Long:             "Open the Cloney documentation in your browser.",
		PersistentPreRun: persistentPreRun,
		RunE:             docsCmdRun,
	}

	return docsCmd
}
