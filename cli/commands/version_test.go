package commands

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testVersionCmd represents a command instance used for testing.
var testVersionCmd = CreateVersionCommand()

// TestVersionCommandOutput tests the output of the version command.
// It should print the version of the app, the OS and the architecture.
func TestVersionCommandOutput(t *testing.T) {
	// Create a new testing.T instance to use with assert functions.
	assert := assert.New(t)

	// Redirect stdout to a buffer.
	var buffer bytes.Buffer
	testVersionCmd.SetOut(&buffer)

	// Execute the command.
	err := testVersionCmd.Execute()
	assert.Nil(err)

	// Assert that the output is correct.
	expectedOutput := fmt.Sprintf("Cloney version %s %s %s\n", appConfig.AppVersion, runtime.GOOS, runtime.GOARCH)
	assert.Equal(expectedOutput, buffer.String())
}
