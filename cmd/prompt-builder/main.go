// ./cmd/prompt-builder/main.go
package main

import (
	"fmt"
	"os"

	"github.com/nnikolov3/prompt-builder/promptbuilder"
)

func main() {
	if err := promptbuilder.RunCLI(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
