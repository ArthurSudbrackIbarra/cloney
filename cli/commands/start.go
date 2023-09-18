package commands

// Start command creates a new cloney template repository.
// Already creates the .cloney.yml file and asks the user
// for informatin such as name, description, etc.
// Flags: path, name, description, authors, license... --yes

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
			fmt.Printf("What is the name of the template repository [%s]: ", appConfig.DefaultCloneyProjectName)
			scanner.Scan()
			name = scanner.Text()
		}

		if description == "" {
			fmt.Printf("What is the description of the template repository [%s]: ", appConfig.DefaultMetadataDescriptionValue)
			scanner.Scan()
			description = scanner.Text()
		}

		if license == "" {
			fmt.Printf("What is the license of the template repository [%s]: ", appConfig.DefaultMetadataLicenseValue)
			scanner.Scan()
			license = scanner.Text()
		}

		if len(authors) == 0 {
			fmt.Print("What are the authors of the template repository (separated by commas): ")
			scanner.Scan()
			authorsStr = scanner.Text()
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

	// Set default values.
	if name == "" {
		name = appConfig.DefaultCloneyProjectName
	}
	if description == "" {
		description = appConfig.DefaultMetadataDescriptionValue
	}
	if license == "" {
		license = appConfig.DefaultMetadataLicenseValue
	}

	// Build the raw YAML metadata string.
	rawMetadata := "# The name of this template repository.\n"
	rawMetadata += fmt.Sprintf("name: %s\n\n", name)
	rawMetadata += "# The description of this template repository.\n"
	rawMetadata += fmt.Sprintf("description: %s\n\n", description)
	rawMetadata += "# The version of this template repository. Change it as you make changes to your template.\n"
	rawMetadata += "template_version: \"0.0.0\"\n\n"
	rawMetadata += "# The version of this manifest. Do not change this value unless you know what you are doing.\n"
	rawMetadata += fmt.Sprintf("manifest_version: %s\n\n", appConfig.MetadataManifestVersion)
	rawMetadata += "# The license of the template repository.\n"
	rawMetadata += fmt.Sprintf("license: %s\n", license)
	if len(authors) > 0 {
		rawMetadata += "\n# The authors of the template repository.\n"
		rawMetadata += "authors:\n"
		for _, author := range authors {
			rawMetadata += fmt.Sprintf("  - %s\n", author)
		}
	}
	// Example variables.
	rawMetadata += "\n# Example variables. Delete this section and add your own variables.\n"
	rawMetadata += "variables:\n"
	rawMetadata += "  - name: app_name\n"
	rawMetadata += "    description: The name of the application.\n"
	rawMetadata += "    default: my-app\n"
	rawMetadata += "    example: my-app # Example is required so that the variable type can be identified.\n\n"
	rawMetadata += "  - name: enable_https\n"
	rawMetadata += "    description: Whether to enable HTTPS or not.\n"
	rawMetadata += "    default: true # Remove default to make the variable required.\n"
	rawMetadata += "    example: true\n"

	// Get the current working directory.
	currentDir, err := steps.GetCurrentWorkingDirectory()
	if err != nil {
		return err
	}

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
	clonePath := steps.CalculateClonePath(repository, currentDir, output)

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
	// Define command-line flags.
	startCmd.Flags().StringP("output", "o", "", "Where to save the template repository")
	startCmd.Flags().StringP("name", "n", "", "The name of the template repository")
	startCmd.Flags().StringP("description", "d", "", "The description of the template repository")
	startCmd.Flags().StringArrayP("authors", "a", []string{}, "The authors of the template repository")
	startCmd.Flags().StringP("license", "l", "", "The license of the template repository")
	startCmd.Flags().BoolP("non-interactive", "y", false, "Skip the questions and use the default values and/or flags")

	// Add the start command to the root command.
	rootCmd.AddCommand(startCmd)
}
