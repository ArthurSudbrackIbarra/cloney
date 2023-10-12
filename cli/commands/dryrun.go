package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/templates"

	"github.com/spf13/cobra"
)

// dryRunCmdRun is the function that runs when the 'dry-run' command is called.
func dryRunCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	var repositorySource string
	if len(args) >= 1 {
		repositorySource = args[0]
	}
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
	sourcePath, _ := steps.CalculatePath(repositorySource, "")
	outputPath, _ := steps.CalculatePath(output, "")

	// Read the repository metadata file.
	metadataFilePath := filepath.Join(sourcePath, appConfig.MetadataFileName)
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

	// Define options for ignoring specific files and directories when filling template variables.
	ignorePaths := []string{
		filepath.Base(filepath.Join(currentDir, variables)),
	}
	ignorePaths = append(ignorePaths, appConfig.KnownIgnorePaths...)
	ignorePaths = append(ignorePaths, cloneyMetadata.Configuration.IgnorePaths...)

	// Check if the output should be displayed in the terminal.
	if outputInTerminal {
		// Fill the template variables and display the output in the terminal instead of creating the files.
		err = steps.FillDirectory(sourcePath, ignorePaths, true, variablesMap)
	} else {
		// Delete the default dry-run output directory if it exists.
		defaultDryRunDir := filepath.Join(sourcePath, appConfig.DefaultDryRunDirectoryName)
		os.RemoveAll(defaultDryRunDir)

		// Create a new directory to save the filled template files.
		err = templates.CopyDirectory(sourcePath, outputPath, ignorePaths)
		if err != nil {
			return fmt.Errorf("error creating output directory %s: %w", outputPath, err)
		}

		// Fill the template variables in the output directory.
		err = steps.FillDirectory(outputPath, ignorePaths, false, variablesMap)

		// Delete files and directories starting with "_" (Ignore Prefix).
		// These are files that should be processed by Cloney but not copied to the output directory.
		steps.DeleteIgnoredPaths(outputPath, ignorePaths)
	}

	if err != nil {
		// If it was not possible to fill the template variables, delete the created directory (if not a dry run).
		if !outputInTerminal {
			os.RemoveAll(outputPath)
		}
		return err
	}

	// Display a completion message if not in terminal output mode.
	if !outputInTerminal {
		cmd.Println("\nDone!")
	}

	return nil
}

// ResetDryRunFlags resets the flags of the 'dry-run' command.
func ResetDryRunFlags(dryRunCmd *cobra.Command) {
	dryRunCmd.Flags().Set("output", appConfig.DefaultDryRunDirectoryName)
	dryRunCmd.Flags().Set("output-in-terminal", "false")
	dryRunCmd.Flags().Set("variables", appConfig.DefaultUserVariablesFileName)
}

// CreateDryRunCommand creates the 'dry-run' command and its respective flags.
func CreateDryRunCommand() *cobra.Command {
	// dryrunCmd represents the dryrun command.
	// It is used to run a template repository in dryrun mode.
	// It is used to check what the output will be with the given variables.
	dryRunCmd := &cobra.Command{
		Use:   "dry-run",
		Short: "Run a template repository in dry-run mode",
		Long: fmt.Sprintf(`Run a template repository in dry-run mode.

The 'cloney dry-run' command is for debugging purposes.
With this command, you can check the output your template repository will generate with the given variables.

By default, 'cloney dry-run' searches for a file named '%s' in your current directory.
You can specify a different file or pass the variables inline as YAML using the '--variables' flag.`, appConfig.DefaultUserVariablesFileName),
		Example: strings.Join([]string{
			"  dry-run",
			"  dry-run ./path/to/my/template",
			"  dry-run ./path/to/my/template -v variables.yaml",
			"  dry-run ./path/to/my/template -v '{ var1: value, var2: value }'",
		}, "\n"),
		Aliases:          []string{"dryrun", "dr", "fill"},
		PersistentPreRun: persistentPreRun,
		RunE:             dryRunCmdRun,
	}

	// Define command-line flags for the 'dryrun' command.
	dryRunCmd.Flags().StringP("output", "o", appConfig.DefaultDryRunDirectoryName, "Path to output the filled template files")
	dryRunCmd.Flags().BoolP("output-in-terminal", "i", false, "Output the filled template file contents in the terminal instead of creating the files")
	dryRunCmd.Flags().StringP("variables", "v", appConfig.DefaultUserVariablesFileName, "Path to a template variables file or raw YAML")

	return dryRunCmd
}
