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

	// If the skip questions flag is not true, ask the user for the information.
	if !nonInteractive {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please answer the following questions to create the template repository:")
		fmt.Print("Press enter to use the default values.\n\n")

		if name == "" {
			fmt.Print("What is the name of the template repository: ")
			scanner.Scan()
			name = scanner.Text()
		}

		if description == "" {
			fmt.Print("What is the description of the template repository: ")
			scanner.Scan()
			description = scanner.Text()
		}

		if license == "" {
			fmt.Print("What is the license of the template repository: ")
			scanner.Scan()
			license = scanner.Text()
		}

		if len(authors) == 0 {
			fmt.Print("What are the authors of the template repository (separated by commas): ")
			scanner.Scan()
			authorsStr = scanner.Text()
			for _, author := range strings.Split(authorsStr, ",") {
				authors = append(authors, strings.TrimSpace(author))
			}
		}

		if err != nil {
			fmt.Println("Error reading user input:", err)
		}
	} else {

	}

	// Get the current working directory.
	currentDir, err := steps.GetCurrentWorkingDirectory()
	if err != nil {
		return err
	}

	// Calculate the directory path.
	outputPath := filepath.Join(currentDir, output)

	fmt.Println(name, description, license, authors, outputPath)

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
	startCmd.Flags().StringP("output", "o", "", "Where to save the template repository.")
	startCmd.Flags().StringP("name", "n", "", "The name of the template repository.")
	startCmd.Flags().StringP("description", "d", "", "The description of the template repository.")
	startCmd.Flags().StringArrayP("authors", "a", []string{}, "The authors of the template repository.")
	startCmd.Flags().StringP("license", "l", "", "The license of the template repository.")
	startCmd.Flags().BoolP("non-interactive", "y", false, "Skip the questions and use the default values and/or flags.")

	// Add the start command to the root command.
	rootCmd.AddCommand(startCmd)
}
