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
			fmt.Printf("What is the name of the template repository [%s]: ", appConfig.DefaultMetadataNameValue)
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
			fmt.Println("Error reading user input:", err)
		}
	}

	// Set default values.
	if name == "" {
		name = appConfig.DefaultMetadataNameValue
	}
	if description == "" {
		description = appConfig.DefaultMetadataDescriptionValue
	}
	if license == "" {
		license = appConfig.DefaultMetadataLicenseValue
	}
	template_version := "0.0.0"
	manifest_version := appConfig.DefaultMetadataManifestVersionValue

	// Build the raw YAML metadata string.
	rawMetadata := "# The name of the template repository.\n"
	rawMetadata += fmt.Sprintf("name: %s\n\n", name)
	rawMetadata += fmt.Sprintf("# The description of the template repository.\n")
	rawMetadata += fmt.Sprintf("description: %s\n\n", description)
	rawMetadata += fmt.Sprintf("# The version of your template repository. Change it as you make changes to your template.\n")
	rawMetadata += fmt.Sprintf("template_version: %s\n\n", template_version)
	rawMetadata += fmt.Sprintf("# The version of this manifest.\n")
	rawMetadata += fmt.Sprintf("manifest_version: %s\n\n", manifest_version)
	rawMetadata += fmt.Sprintf("# The license of the template repository.\n")
	rawMetadata += fmt.Sprintf("license: %s\n", license)
	if len(authors) > 0 {
		rawMetadata += "authors:\n"
		for _, author := range authors {
			rawMetadata += fmt.Sprintf("  - %s\n", author)
		}
	}
	// Example variables.
	rawMetadata += "\n# Example variables.\n"
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

	// Calculate the directory paths.
	var outputPath string
	if output == "" {
		outputPath = filepath.Join(currentDir, appConfig.DefaultCloneyProjectName)
	} else {
		outputPath = filepath.Join(currentDir, output)
	}
	metadataFilePath := filepath.Join(outputPath, appConfig.MetadataFileName)

	// Create the template repository directory.
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating template repository directory:", err)
		return err
	}
	err = os.WriteFile(metadataFilePath, []byte(rawMetadata), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating metadata file:", err)
		return err
	}

	fmt.Println("\nDone!")

	return nil
}

// startCmd represents the start command.
// This command is used to create a new cloney template repository.
var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Creates a new cloney template repository",
	Long:    "\nCreates a new cloney template repository.",
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
