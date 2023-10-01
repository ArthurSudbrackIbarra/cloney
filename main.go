package main

import (
	"fmt"
	"os"

	"github.com/ArthurSudbrackIbarra/cloney/cli"
	"github.com/ArthurSudbrackIbarra/cloney/config"
)

// main is the entry point of the application.
func main() {
	// Load the application configuration.
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading application configuration: %v\n", err)
		os.Exit(1)
	}

	// Start the CLI.
	cli.Initialize()
}
