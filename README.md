# Prompt Builder

A standalone Go application for building structured prompts from various components including files, system messages, guidelines, and task presets. This tool is designed to be reusable in larger AI helper applications.

## Features

- **File Integration**: Include file contents in prompts with automatic code fencing
- **System Presets**: Predefined system messages for common tasks (coding, analysis, documentation)
- **Custom Guidelines**: Add specific instructions and constraints
- **Multiple Output Formats**: JSON, text, and markdown output
- **Security**: File content fencing and validation with path traversal protection
- **Extensible**: Easy to add new system presets and file types
- **Code Quality**: Comprehensive linting and testing with golangci-lint

## Installation

### Prerequisites

- Go 1.21 or later
- git

### Build from Source

```bash
git clone <repository-url>
cd prompt-builder
make build
```

Or manually:

```bash
git clone <repository-url>
cd prompt-builder
go build -o prompt-builder ./cmd/prompt-builder
```

### Install to System

```bash
make install
```

This will install the binary to `~/bin/prompt-builder`.

## Usage

### Basic Usage

```bash
# Simple prompt
./prompt-builder -p "Explain this code"

# With file content
./prompt-builder -p "Explain this code" -f main.go

# With task preset
./prompt-builder -p "Write a function" -t coding

# With custom system message
./prompt-builder -p "Analyze this" -sys "You are an expert analyst"
```

### Advanced Usage

```bash
# Complete example with all features
./prompt-builder \
  -p "Refactor this code" \
  -f app.py \
  -t coding \
  -g "Follow PEP 8 and add comprehensive comments" \
  -o json
```

### Command Line Options

| Flag | Long Flag | Description |
|------|-----------|-------------|
| `-p` | `--prompt` | User prompt text (required) |
| `-f` | `--file` | Optional file to include in context |
| `-t` | `--task` | Task preset for system message |
| `-sys` | `--system` | Custom system message |
| `-g` | `--guidelines` | Guidelines to follow |
| `-o` | `--output` | Output format (json, text, markdown) |
| `-h` | `--help` | Show help message |

### Task Presets

The following system presets are available:

- **coding**: "You are an expert software developer. Write clean, efficient, and well-documented code."
- **analysis**: "You are an expert code analyst. Provide detailed analysis and insights."
- **documentation**: "You are an expert technical writer. Create clear and comprehensive documentation."

### Output Formats

- **json**: Structured JSON output with all prompt components
- **text**: Plain text format with the assembled prompt
- **markdown**: Markdown format with code fencing (default)

## Architecture

### Core Components

- **CLIFlags**: Command line interface flag parsing and validation
- **BuildRequest**: Internal request structure for building prompts
- **PromptBuilder**: Main orchestrator for assembling prompts
- **FileProcessor**: Handles file reading, validation, and content fencing
- **Prompt**: Final assembled prompt with all components

### File Structure

```
prompt-builder/
├── cmd/
│   └── prompt-builder/
│       └── main.go              # CLI entry point
├── internal/
│   └── prompt-builder/
│       ├── types.go             # Core data structures
│       ├── builder.go           # Prompt building logic
│       ├── file_processor.go    # File handling and security
│       ├── cli.go               # CLI functionality
│       ├── cli_test.go          # CLI tests
│       ├── file_processor_test.go # File processor tests
│       └── types_test.go        # Type validation tests
├── .golangci.yml                # Linting configuration
├── Makefile                     # Build and development tasks
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Security Features

The application includes several security measures:

- **Path Traversal Protection**: Prevents access to system directories
- **File Size Limits**: Configurable maximum file size (default: 1MB)
- **Extension Whitelist**: Only allows specific file extensions
- **Content Fencing**: Wraps file content with BEGIN/END markers
- **Home Directory Restriction**: Limits file access to user's home directory

## Integration with AI Helpers

This prompt-builder is designed to be easily integrated into larger AI helper projects. The internal package can be imported and used as a library:

```go
import "github.com/nnikolov3/prompt-builder/internal/prompt-builder"

// Create a file processor with custom settings
fileProcessor := promptbuilder.NewFileProcessor(
    1024*1024, // 1MB max file size
    []string{".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h", ".cs", ".php", ".rb", ".rs", ".txt", ".md"},
)

// Create a prompt builder
builder := promptbuilder.NewPromptBuilder(fileProcessor)

// Add custom system presets
builder.AddSystemPreset("custom-task", "Custom system message")

// Build a prompt
req := &promptbuilder.BuildRequest{
    Prompt:        "Your prompt here",
    File:          "path/to/file.go",
    Task:          "coding",
    SystemMessage: "Custom system message",
    Guidelines:    "Follow best practices",
    OutputFormat:  "json",
}

result, err := builder.BuildPrompt(req)
if err != nil {
    // Handle error
}

// Use the assembled prompt
promptText := result.Prompt.String()
```

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for code quality checks)

### Available Make Targets

```bash
# Build the application
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Format code
make fmt

# Clean build artifacts
make clean

# Install to ~/bin
make install

# Show all available targets
make help
```

### Testing

Run all tests:

```bash
make test
```

Or manually:

```bash
go test -v ./...
```

Run specific test suites:

```bash
# CLI tests
go test -v -run TestCLI ./internal/prompt-builder

# File processor tests
go test -v -run TestFile ./internal/prompt-builder

# All validation tests
go test -v -run Test.*Validate ./internal/prompt-builder
```

### Code Quality

The project uses golangci-lint for comprehensive code quality checks:

```bash
make lint
```

This runs all configured linters including:
- cyclop (cyclomatic complexity)
- err113 (error handling)
- exhaustruct (struct initialization)
- forbidigo (forbidden functions)
- funlen (function length)
- gocognit (cognitive complexity)
- lll (line length)
- tagliatelle (JSON tag formatting)
- varnamelen (variable naming)
- wsl_v5 (whitespace and formatting)

### Adding New Features

1. **Write tests first** following TDD principles
2. **Implement the feature** in the smallest possible module
3. **Run tests** to ensure they pass: `make test`
4. **Lint the code** to ensure quality: `make lint`
5. **Refactor** for clarity and simplicity
6. **Repeat** the cycle until the code meets all quality standards

### Design Principles

This project follows strict design principles:

- **Simplicity**: Prefer the simplest solution that meets requirements
- **Explicitness**: Make assumptions visible and auditable
- **Modularity**: Isolate responsibilities into small, composable units
- **Correctness**: Write tests before or alongside code
- **Maintainability**: Keep functions short, focused, and intention-revealing
- **Consistency**: Follow established patterns and conventions
- **Security**: Validate all inputs and prevent common vulnerabilities

### Code Quality Standards

- All code must pass `make lint` (golangci-lint)
- All tests must pass
- No magic numbers or hardcoded values
- Self-documenting code with clear comments
- No unused imports or dead code
- Proper error handling with wrapped static errors
- Consistent JSON tag formatting (camelCase)
- Descriptive variable names

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
