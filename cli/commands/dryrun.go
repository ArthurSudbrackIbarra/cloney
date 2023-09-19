package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/templates"
	"github.com/ArthurSudbrackIbarra/cloney/utils"
	"github.com/spf13/cobra"
)

// dryrunCmdRun is the function that runs when the dryrun command is called.
func dryrunCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	path, _ := cmd.Flags().GetString("path")
	output, _ := cmd.Flags().GetString("output")
	outputInTerminal, _ := cmd.Flags().GetBool("output-in-terminal")
	variables, _ := cmd.Flags().GetString("variables")

	// Variable to store errors.
	var err error

	// Get the current working directory.
	currentDir, err := steps.GetCurrentWorkingDirectory()
	if err != nil {
		return err
	}

	// Get the template variables provided by the user.
	variablesMap, err := steps.GetUserVariablesMap(currentDir, variables)
	if err != nil {
		return err
	}

	// Calculate the directory paths.
	targetPath := filepath.Join(currentDir, path)
	outputPath := filepath.Join(currentDir, output)

	// Read the repository metadata file.
	metadataFilePath := filepath.Join(targetPath, appConfig.MetadataFileName)
	metadataContent, err := steps.ReadRepositoryMetadata(metadataFilePath)
	if err != nil {
		return err
	}

	// Parse the metadata file.
	cloneyMetadata, err := steps.ParseRepositoryMetadata(metadataContent, appConfig.SupportedManifestVersions)
	if err != nil {
		return err
	}

	// Validate if the user variables match the template variables.
	// Also, fill default values of the variables if they are not defined.
	err = steps.MatchUserVariables(cloneyMetadata, variablesMap)
	if err != nil {
		return err
	}

	// Fill the template variables.
	templateOptions := templates.TemplateFillOptions{
		SourceDirectoryPath: targetPath,
		TargetDirectoryPath: &outputPath,
		TerminalMode:        outputInTerminal,
	}
	ignoreOptions := utils.IgnorePathOptions{
		// Ignore specific files when filling the template variables.
		IgnoreFiles: []string{
			appConfig.MetadataFileName,
			appConfig.DefaultUserVariablesFileName,
			filepath.Base(filepath.Join(currentDir, variables)),
		},
		// Ignore '.git' directory when filling the template variables.
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

	if !outputInTerminal {
		fmt.Println("\nDone!")
	}

	return nil
}

// dryrunCmd represents the dryrun command.
// It is used to run a template repository in dryrun mode.
// It is used to check what the output will be with the given variables.
var dryrunCmd = &cobra.Command{
	Use:   "dryrun",
	Short: "Run a template repository in dryrun mode",
	Long: fmt.Sprintf(`Run a template repository in dryrun mode.

The 'cloney dryrun' command is for debugging purposes.
With this command, you can check the output your template repository will generate with the given variables.

By default, 'cloney dryrun' searches for a file named '%s' in your current directory.
You can specify a different file using the '--variables' flag or pass the variables inline as YAML.`, appConfig.DefaultUserVariablesFileName),
	Example: " cloney dryrun",
	Aliases: []string{"dry-run", "dr"},
	RunE:    dryrunCmdRun,
}

// InitializeDryrun initializes the dryrun command.
func InitializeDryrun(rootCmd *cobra.Command) {
	appConfig := config.GetAppConfig()
	// Define command-line flags for the 'dryrun' command.
	dryrunCmd.Flags().StringP("path", "p", "", "Path to your local template repository")
	dryrunCmd.Flags().StringP("output", "o", appConfig.DefaultDryRunDirectoryName, "Path to output the filled template files")
	dryrunCmd.Flags().BoolP("output-in-terminal", "i", false, "Output the filled template file contents in the terminal instead of creating the files")
	dryrunCmd.Flags().StringP("variables", "v", appConfig.DefaultUserVariablesFileName, "Path to a template variables file or raw YAML")

	// Add the 'dryrun' command to the root command.
	rootCmd.AddCommand(dryrunCmd)
}
