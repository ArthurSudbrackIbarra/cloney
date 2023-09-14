package commands

import (
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/templates"
	"github.com/ArthurSudbrackIbarra/cloney/utils"
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
	currentDir, err := steps.GetCurrentWorkingDirectory()
	if err != nil {
		return err
	}

	// Get the template variables provided by the user.
	variablesMap, err := steps.GetUserVariablesMap(currentDir, variablesJSON, variablesFilePath)
	if err != nil {
		return err
	}

	// Calculate the directory paths.
	targetPath := filepath.Join(currentDir, path)
	outputPath := filepath.Join(currentDir, output)

	// Read the repository metadata file.
	appConfig := config.GetAppConfig()
	metadataFilePath := filepath.Join(targetPath, appConfig.MetadataFileName)
	metadataContent, err := steps.ReadRepositoryMetadata(metadataFilePath)
	if err != nil {
		return err
	}

	// Parse the metadata file.
	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent)
	if err != nil {
		return err
	}

	// Validate if the user variables match the template variables.
	// Also fill default values of the variables if they are not defined.
	err = steps.MatchUserVariables(cloneyMetadata, variablesMap)
	if err != nil {
		return err
	}

	// Fill the template variables.
	// If ouput in terminal is enabled, the filled template files will be printed to the terminal instead of being saved to the files.
	templateOptions := templates.TemplateFillOptions{
		SourceDirectoryPath: targetPath,
		TargetDirectoryPath: &outputPath,
		TerminalMode:        outputInTerminal,
	}
	ignoreOptions := utils.IgnorePathOptions{
		// Ignore the metadata file when filling the template variables.
		// Ignore the user variables file when filling the template variables.
		IgnoreFiles: []string{appConfig.MetadataFileName, filepath.Base(variablesFilePath)},

		// Ignore '.git' directories when filling the template variables.
		IgnoreDirectories: []string{".git"},
	}
	err = steps.FillTemplateVariables(templateOptions, ignoreOptions, variablesMap)
	if err != nil {
		// If it was not possible to fill the template variables, delete the created directory.
		if !outputInTerminal {
			os.RemoveAll(outputPath)
		}
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
	Aliases: []string{"dry-run", "dr"},
	RunE:    dryrunCmdRun,
}

// InitializeDryrun initializes the dryrun command.
func InitializeDryrun(rootCmd *cobra.Command) {
	appConfig := config.GetAppConfig()
	// Define command-line flags.
	dryrunCmd.Flags().StringP("path", "p", "", "Path to your local template repository")
	dryrunCmd.Flags().StringP("output", "o", "dry-run-result", "Path to output the filled template files")
	dryrunCmd.Flags().BoolP("output-in-terminal", "i", false, "Output the filled template file contents in the terminal instead of creating the files")
	dryrunCmd.Flags().StringP("variables-file", "f", appConfig.DefaultUserVariablesFileName, "Path to a template variables YAML or JSON file")
	dryrunCmd.Flags().StringP("variables", "v", "", "Inline template variables as JSON")

	// Add command to the root command.
	rootCmd.AddCommand(dryrunCmd)
}
