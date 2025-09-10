// ./cmd/prompt-builder/main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nnikolov3/prompt-builder/internal/promptbuilder"
)

func main() {
	// 1. Define and parse command-line flags.
	templateStr := flag.String("t", "", "The prompt template string.")
	flag.Parse()

	if *templateStr == "" {
		fmt.Fprintln(os.Stderr, "Error: template string is required. Use the -t flag.")
		os.Exit(1)
	}

	// 2. Read variables from standard input.
	stdinBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	var vars map[string]any
	if err := json.Unmarshal(stdinBytes, &vars); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON from stdin: %v\n", err)
		os.Exit(1)
	}

	// 3. Use the library to do the work.
	builder, err := promptbuilder.New(*templateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating builder: %v\n", err)
		os.Exit(1)
	}
	builder.SetMap(vars)

	result, err := builder.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building prompt: %v\n", err)
		os.Exit(1)
	}

	// 4. Print the final result to standard output.
	fmt.Print(result)
}
