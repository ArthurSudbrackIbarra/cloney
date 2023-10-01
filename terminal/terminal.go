package terminal

import (
	"bufio"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Colors for terminal output.
var (
	Green  = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
)

// cmd is the current command being executed.
var cmd *cobra.Command

// SetCmd sets the current command, allowing messages to be printed to the command's output using cmd.Print().
func SetCmd(newCommand *cobra.Command) {
	cmd = newCommand
}

// Messages functions.

// Message prints a message with no prefix.
func Message(message string) {
	if cmd != nil {
		cmd.Print(fmt.Sprintf("%s\n", message))
	}
}

// OKMessage prints a success message with a green "[OK]" prefix.
func OKMessage(message string) {
	if cmd != nil {
		cmd.Print(fmt.Sprintf("[%s] %s\n", Green("OK"), message))
	}
}

// WarningMessage prints a warning message with a yellow "[Warning]" prefix.
func WarningMessage(message string) {
	if cmd != nil {
		cmd.Print(fmt.Sprintf("[%s] %s\n", Yellow("Warning"), message))
	}
}

// ErrorMessage prints an error message with a red "[Error]" prefix.
func ErrorMessage(message string, err error) {
	if cmd != nil {
		cmd.Print(fmt.Sprintf("[%s] %s: %v\n", Red("Error"), message, err))
	}
}

// InputWithDefaultValue prompts the user for input via terminal and returns the input value or the default value.
func InputWithDefaultValue(scanner *bufio.Scanner, message, defaultValue string) string {
	if cmd != nil {
		cmd.Print(fmt.Sprintf("%s [%s]: ", message, Blue(defaultValue)))
	}
	scanner.Scan()
	input := scanner.Text()
	if input == "" {
		return defaultValue
	}
	return input
}
