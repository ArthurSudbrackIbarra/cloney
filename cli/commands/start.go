package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArthurSudbrackIbarra/cloney/cli/commands/steps"
	"github.com/ArthurSudbrackIbarra/cloney/config"
	"github.com/ArthurSudbrackIbarra/cloney/utils"
	"github.com/spf13/cobra"
)

// startCmdRun is the function that runs when the start command is called.
func startCmdRun(cmd *cobra.Command, args []string) error {
	// Get command-line arguments.
	output, _ := cmd.Flags().GetString("output")
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	var authorsStr string
	authors, _ := cmd.Flags().GetStringArray("authors")
	license, _ := cmd.Flags().GetString("license")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Variable to store errors.
	var err error

	// If the non-interactive flag is not set, ask the user for the information.
	appConfig := config.GetAppConfig()
	if !nonInteractive {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please answer the following questions to create the template repository.")
		fmt.Print("Press enter to use the default values.\n\n")

		if name == "" {
			name = steps.InputWithDefaultValue(
				scanner, "What is the name of the template repository", appConfig.DefaultCloneyProjectName,
			)
		}

		if description == "" {
			description = steps.InputWithDefaultValue(
				scanner, "What is the description of the template repository", appConfig.DefaultMetadataDescriptionValue,
			)
		}

		if license == "" {
			license = steps.InputWithDefaultValue(
				scanner, "What is the license of the template repository", appConfig.DefaultMetadataLicenseValue,
			)
		}

		if len(authors) == 0 {
			authorsStr = steps.InputWithDefaultValue(scanner, "What are the authors of the template repository (separated by commas)", "")
			if authorsStr != "" {
				for _, author := range strings.Split(authorsStr, ",") {
					authors = append(authors, strings.TrimSpace(author))
				}
			}
		}

		if err != nil {
			utils.ErrorMessage("Error reading user input", err)
		}
	}

	// Build the raw YAML metadata string.
	rawMetadata := "# The version of this Cloney manifest file, ensuring compatibility with different versions of Cloney.\n"
	rawMetadata += fmt.Sprintf("manifest_version: %s\n\n", appConfig.MetadataManifestVersion)
	rawMetadata += "# The name of your template, providing a clear identifier for users.\n"
	rawMetadata += fmt.Sprintf("name: %s\n\n", name)
	rawMetadata += "# A brief but informative description of your template's purpose and functionality.\n"
	rawMetadata += fmt.Sprintf("description: %s\n\n", description)
	rawMetadata += "# The version number of your template. Update it as you make new changes to your template.\n"
	rawMetadata += "template_version: \"0.0.0\"\n\n"
	rawMetadata += "# The licensing information for your template, specifying how others can use and distribute it.\n"
	rawMetadata += fmt.Sprintf("license: %s\n", license)
	if len(authors) > 0 {
		rawMetadata += "\n# A list of contributors or creators of the template, acknowledging their role in its development.\n"
		rawMetadata += "authors:\n"
		for _, author := range authors {
			rawMetadata += fmt.Sprintf("  - %s\n", author)
		}
	}
	// Example variables.
	rawMetadata += "\n# A list of variables that users can customize during the cloning process.\n"
	rawMetadata += "# Delete this section and add your own variables.\n"
	rawMetadata += "variables:\n"
	rawMetadata += "  - name: app_name\n"
	rawMetadata += "    description: The name of the application.\n"
	rawMetadata += "    default: my-app\n"
	rawMetadata += "    example: my-app\n\n"
	rawMetadata += "  - name: enable_https\n"
	rawMetadata += "    description: Whether to enable HTTPS or not.\n"
	rawMetadata += "    example: true\n"

	// Suppress prints for common-steps functions.
	steps.SetSuppressPrints(true)

	// Create and validate a reference to the Cloney example repository.
	repository, err := steps.CreateAndValidateRepository(
		appConfig.CloneyExampleRepositoryURL, "main", "",
	)
	if err != nil {
		fmt.Println("Error when cloning the example Cloney repository from GitHub.")
		return err
	}

	// Calculate the clone path.
	if output == "" {
		// If the output flag is not set, use the name of the template repository as the name of the directory.
		output = name
	}
	clonePath, _ := steps.CalculatePath(output, repository.GetName())

	// Clone the repository.
	err = steps.CloneRepository(repository, clonePath)
	if err != nil {
		return err
	}

	// Delete the .git directory.
	gitDirPath := filepath.Join(clonePath, ".git")
	os.RemoveAll(gitDirPath)

	// Update the metadata file.
	metadataFilePath := filepath.Join(clonePath, appConfig.MetadataFileName)
	err = os.WriteFile(metadataFilePath, []byte(rawMetadata), os.ModePerm)
	if err != nil {
		utils.ErrorMessage("Error creating the repository metadata file", err)
		return err
	}

	fmt.Println("\nDone!")

	return nil
}

// startCmd represents the start command.
// This command is used to create a new cloney template repository.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Creates a new cloney template repository",
	Long: `Creates a new cloney template repository.

cloney start will create a directory with the necessary files to start a new cloney template repository.`,
	Example: "  cloney start",
	RunE:    startCmdRun,
}

// InitializeStart initializes the start command.
func InitializeStart(rootCmd *cobra.Command) {
	// Define command-line flags for the 'start' command.
	startCmd.Flags().StringP("output", "o", "", "Where to save the template repository")
	startCmd.Flags().StringP("name", "n", "", "The name of the template repository")
	startCmd.Flags().StringP("description", "d", "", "The description of the template repository")
	startCmd.Flags().StringArrayP("authors", "a", []string{}, "The authors of the template repository")
	startCmd.Flags().StringP("license", "l", "", "The license of the template repository")
	startCmd.Flags().BoolP("non-interactive", "y", false, "Skip the questions and use the default values and/or flags")

	// Add the 'start' command to the root command.
	rootCmd.AddCommand(startCmd)
}
