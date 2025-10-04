// ./cmd/prompt-builder/main.go
package main

import (
	"fmt"
	"os"

	"github.com/book-expert/prompt-builder/promptbuilder"
)

func main() {
	err := promptbuilder.RunCLI(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
