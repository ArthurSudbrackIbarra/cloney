package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"
	"github.com/ArthurSudbrackIbarra/cloney/templates"
	"github.com/spf13/cobra"
)

// Dryrun command is used to run a template repository in dryrun mode.
// Check what the output will be with the given variables.
// Flags: path, var-..., var-file, output.

// dryrunCmdRun is the function that runs when the dryrun command is called.
func dryrunCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	path, _ := cmd.Flags().GetString("path")
	output, _ := cmd.Flags().GetString("output")
	outputInTerminal, _ := cmd.Flags().GetBool("output-in-terminal")
	variablesFilePath, _ := cmd.Flags().GetString("variables-file")
	variablesJSON, _ := cmd.Flags().GetString("variables")

	// Variable to store errors.
	var err error

	// Get the current working directory.
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get user's current directory:", err)
		return err
	}

	// Get the template variables provided by the user.
	var variablesMap map[string]interface{}
	if variablesJSON != "" {
		variablesMap, err = metadata.NewCloneyUserVariablesFromRawJSON(variablesJSON)
		if err != nil {
			fmt.Println("Could not parse template variables:", err)
			return err
		}
	} else {
		variablesFilePath = filepath.Join(currentDir, variablesFilePath)
		variablesMap, err = metadata.NewCloneyUserVariablesFromFile(variablesFilePath)
		if err != nil {
			fmt.Println("Could not read template variables file:", err)
			return err
		}
	}

	// Calculate the directory paths.
	targetPath := filepath.Join(currentDir, path)
	outputPath := filepath.Join(currentDir, output)

	// Fill the template variables.
	// If only-in-terminal is enabled, the filled template files will be printed to the terminal instead of being saved to the files.
	filler := templates.NewTemplateFiller(variablesMap)
	options := templates.TemplateFillOptions{
		SourceDirectoryPath: targetPath,
		TargetDirectoryPath: &outputPath,
		TerminalMode:        outputInTerminal,
	}
	err = filler.FillDirectory(options)
	if err != nil {
		fmt.Println("Could not fill template:", err)
		return err
	}

	return nil
}

// dryrunCmd represents the dryrun command.
// It is used to run a template repository in dryrun mode.
// It is used to check what the output will be with the given variables.
var dryrunCmd = &cobra.Command{
	Use:     "dryrun",
	Short:   "Run a template repository in dryrun mode",
	Long:    "\nRun a template repository in dryrun mode.",
	Example: " cloney dryrun",
	RunE:    dryrunCmdRun,
}

// InitializeDryrun initializes the dryrun command.
func InitializeDryrun(rootCmd *cobra.Command) {
	appConfig := config.GetAppConfig()
	// Define command-line flags.
	dryrunCmd.Flags().StringP("path", "p", "", "Path to your local template repository")
	dryrunCmd.Flags().StringP("output", "o", "", "Path to output the filled template files")
	dryrunCmd.Flags().BoolP("output-in-terminal", "i", false, "Output the filled template file contents in the terminal instead of creating the files")
	dryrunCmd.Flags().StringP("variables-file", "f", appConfig.DefaultUserVariablesFileName, "Path to a template variables YAML or JSON file")
	dryrunCmd.Flags().StringP("variables", "v", "", "Inline template variables as JSON")

	// Add command to the root command.
	rootCmd.AddCommand(dryrunCmd)
}
