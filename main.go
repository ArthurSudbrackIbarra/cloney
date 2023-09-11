package main

import (
	"github.com/ArthurSudbrackIbarra/cloney/cli"
	"github.com/ArthurSudbrackIbarra/cloney/config"
)

// main is the entry point of the application.
func main() {
	// Load the application configuration.
	err := config.LoadConfig()
	if err != nil {
		panic("Could not load application configuration.")
	}

	// Start the CLI.
	cli.Initialize()
}
