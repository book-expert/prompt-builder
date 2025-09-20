# Prompt Builder

## Project Summary

A standalone Go application and library for building structured prompts from various components.

## Detailed Description

This tool is designed to be a reusable component in larger AI helper applications. It allows for the construction of detailed and structured prompts by combining various inputs such as user prompts, file contents, system messages, and guidelines.

Key features include:

-   **File Integration**: Include file contents in prompts with automatic code fencing.
-   **System Presets**: Predefined system messages for common tasks (e.g., coding, analysis, documentation).
-   **Custom Guidelines**: Add specific instructions and constraints to the prompt.
-   **Multiple Output Formats**: Supports JSON, text, and markdown output.
-   **Security**: Includes file content fencing and validation with path traversal protection.

## Technology Stack

-   **Programming Language:** Go 1.25

## Getting Started

### Prerequisites

-   Go 1.25 or later.

### Installation

To build the application, you can use the `make build` command:

```bash
make build
```

This will create the `prompt-builder` binary in your `~/bin` directory.

To use this library in your project, you can use `go get`:

```bash
go get github.com/book-expert/prompt-builder
```

## Usage

### Command-Line Interface

```bash
# Simple prompt
prompt-builder -p "Explain this code"

# With file content
prompt-builder -p "Explain this code" -f main.go

# With task preset
prompt-builder -p "Write a function" -t coding

# With custom system message
prompt-builder -p "Analyze this" -sys "You are an expert analyst"
```

### Library

```go
package main

import (
    "fmt"

    "github.com/book-expert/prompt-builder/promptbuilder"
)

func main() {
    fileProcessor := promptbuilder.NewFileProcessor(
        1024*1024, // 1MB max file size
        []string{".go", ".py", ".js"},
    )

    builder := promptbuilder.New(fileProcessor)
    builder.AddSystemPreset("coding", "You are an expert software developer.")

    req := &promptbuilder.BuildRequest{
        Prompt:     "Refactor this code",
        File:       "app.py",
        Task:       "coding",
        Guidelines: "Follow PEP 8",
    }

    result, err := builder.BuildPrompt(req)
    if err != nil {
        panic(err)
    }

    fmt.Println(result.Prompt.String())
}
```

## Testing

To run the tests for this library, you can use the `make test` command:

```bash
make test
```

## License

Distributed under the MIT License. See the `LICENSE` file for more information.