package terminal

import (
	"bufio"
	"bytes"
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

	BlueBoldUnderline  = color.New(color.FgBlue, color.Bold, color.Underline).SprintFunc()
	WhiteBoldUnderline = color.New(color.FgWhite, color.Bold, color.Underline).SprintFunc()
)

// cmd is the current command being executed.
var cmd *cobra.Command

// testBuffer is a buffer used for testing.
var testBuffer *bytes.Buffer

// SetCmd sets the current command, allowing messages to be printed to the command's output using cmd.Print().
func SetCmd(newCommand *cobra.Command) {
	cmd = newCommand
}

// SetTestMode sets the terminal to test mode, allowing messages to be printed to a buffer.
func SetTestMode(newTestBuffer *bytes.Buffer) {
	testBuffer = newTestBuffer
}

// Messages functions.

// Message prints a message with no prefix.
func Message(message string) {
	if cmd != nil {
		str := fmt.Sprintf("%s\n", message)

		// Print the message to the terminal.
		cmd.SetOut(cmd.OutOrStdout())
		cmd.Print(str)

		// If in test mode, write to the buffer as well.
		if testBuffer != nil {
			cmd.SetOut(testBuffer)
			cmd.Print(str)
		}
	}
}

// Messagef prints a formatted message with no prefix.
func Messagef(format string, a ...interface{}) {
	if cmd != nil {
		str := fmt.Sprintf(format, a...)

		// Print the message to the terminal.
		cmd.SetOut(cmd.OutOrStdout())
		cmd.Print(str)

		// If in test mode, write to the buffer as well.
		if testBuffer != nil {
			cmd.SetOut(testBuffer)
			cmd.Print(str)
		}
	}
}

// OKMessage prints a success message with a green "[OK]" prefix.
func OKMessage(message string) {
	if cmd != nil {
		str := fmt.Sprintf("[%s] %s\n", Green("OK"), message)

		// Print the message to the terminal.
		cmd.SetOut(cmd.OutOrStdout())
		cmd.Print(str)

		// If in test mode, write to the buffer as well.
		if testBuffer != nil {
			cmd.SetOut(testBuffer)
			cmd.Print(str)
		}
	}
}

// WarningMessage prints a warning message with a yellow "[Warning]" prefix.
func WarningMessage(message string) {
	if cmd != nil {
		str := fmt.Sprintf("[%s] %s\n", Yellow("Warning"), message)

		// Print the message to the terminal.
		cmd.SetOut(cmd.OutOrStdout())
		cmd.Print(str)

		// If in test mode, write to the buffer as well.
		if testBuffer != nil {
			cmd.SetOut(testBuffer)
			cmd.Print(str)
		}
	}
}

// ErrorMessage prints an error message with a red "[Error]" prefix.
func ErrorMessage(message string, err error) {
	if cmd != nil {
		str := ""
		if err != nil {
			str = fmt.Sprintf("[%s] %s: %v\n", Red("Error"), message, err)
		} else {
			str = fmt.Sprintf("[%s] %s\n", Red("Error"), message)
		}

		// Print the message to the terminal.
		cmd.SetOut(cmd.ErrOrStderr())
		cmd.Print(str)

		// If in test mode, write to the buffer as well.
		if testBuffer != nil {
			cmd.SetOut(testBuffer)
			cmd.Print(str)
		}
	}
}

// CautionInput prompts the user for input via terminal with a caution prefix and returns the input.
func CautionInput(scanner *bufio.Scanner, message string) string {
	if cmd != nil {
		cmd.SetOut(cmd.OutOrStdout())
		cmd.Print(fmt.Sprintf("[%s] %s: ", Yellow("Caution"), message))
	}
	scanner.Scan()
	return scanner.Text()
}

// InputWithDefaultValue prompts the user for input via terminal and returns the input value or the default value.
func InputWithDefaultValue(scanner *bufio.Scanner, message, defaultValue string) string {
	if cmd != nil {
		cmd.SetOut(cmd.OutOrStdout())
		cmd.Print(fmt.Sprintf("%s [%s]: ", message, Blue(defaultValue)))
	}
	scanner.Scan()
	input := scanner.Text()
	if input == "" {
		return defaultValue
	}
	return input
}
