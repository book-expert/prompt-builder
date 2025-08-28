// Package main provides the entry point for the prompt-builder CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/nnikolov3/prompt-builder/internal/prompt-builder"
)

func main() {
	// Get command line arguments (skip the program name)
	args := os.Args[1:]

	// Run the CLI application
	err := promptbuilder.RunCLI(args, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
