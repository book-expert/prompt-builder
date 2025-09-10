package promptbuilder

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
)

// ParseFlags parses command line arguments into CLIFlags.
func ParseFlags(args []string) (*CLIFlags, error) {
	flagSet := flag.NewFlagSet("prompt-builder", flag.ExitOnError)

	var flags CLIFlags

	flagSet.StringVar(&flags.Prompt, "p", "", "User prompt text (required)")
	flagSet.StringVar(&flags.Prompt, "prompt", "", "User prompt text (required)")
	flagSet.StringVar(&flags.File, "f", "", "Optional file to include in context")
	flagSet.StringVar(&flags.File, "file", "", "Optional file to include in context")
	flagSet.StringVar(&flags.Task, "t", "", "Task preset for system message")
	flagSet.StringVar(&flags.Task, "task", "", "Task preset for system message")
	flagSet.StringVar(&flags.SystemMessage, "sys", "", "Custom system message")
	flagSet.StringVar(&flags.SystemMessage, "system", "", "Custom system message")
	flagSet.StringVar(&flags.Guidelines, "g", "", "Guidelines to follow")
	flagSet.StringVar(&flags.Guidelines, "guidelines", "", "Guidelines to follow")
	flagSet.StringVar(&flags.OutputFormat, "o", "", "Output format (json, text, markdown)")
	flagSet.StringVar(&flags.OutputFormat, "output", "", "Output format (json, text, markdown)")

	// Parse the flags
	err := flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Validate the flags
	if err := flags.Validate(); err != nil {
		return nil, fmt.Errorf("invalid flags: %w", err)
	}

	return &flags, nil
}

// PrintUsage prints the usage information for the CLI.
func PrintUsage() {
	log.Print(`Usage: prompt-builder [OPTIONS]

Build prompts from various components including files, system messages, and guidelines.

OPTIONS:
  -p, --prompt TEXT          User prompt text (required)
  -f, --file PATH           Optional file to include in context
  -t, --task TASK           Task preset for system message
  -sys, --system TEXT       Custom system message
  -g, --guidelines TEXT     Guidelines to follow
  -o, --output FORMAT       Output format (json, text, markdown)
  -h, --help                Show this help message

EXAMPLES:
  prompt-builder -p "Explain this code" -f main.go
  prompt-builder -p "Refactor this" -f app.py -t coding -g "Follow PEP 8"
  prompt-builder -p "Analyze this code" -f app.js -o json
`)
}

// RunCLI runs the CLI application with the given arguments and writes output to the provided writer.
func RunCLI(args []string, output io.Writer) error {
	// Check for help flag
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			PrintUsage()

			return nil
		}
	}

	// Parse flags
	flags, err := ParseFlags(args)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Create file processor with reasonable defaults
	allowedExtensions := []string{
		".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h", ".cs", ".php", ".rb", ".rs", ".txt", ".md",
	}
	fileProcessor := NewFileProcessor(1024*1024, allowedExtensions)

	// Create prompt builder
	builder := NewPromptBuilder(fileProcessor)

	// Add some default system presets
	codingPreset := "You are an expert software developer. Write clean, efficient, and well-documented code."
	if err := builder.AddSystemPreset("coding", codingPreset); err != nil {
		return fmt.Errorf("failed to add coding preset: %w", err)
	}

	analysisPreset := "You are an expert code analyst. Provide detailed analysis and insights."
	if err := builder.AddSystemPreset("analysis", analysisPreset); err != nil {
		return fmt.Errorf("failed to add analysis preset: %w", err)
	}

	documentationPreset := "You are an expert technical writer. Create clear and comprehensive documentation."
	if err := builder.AddSystemPreset("documentation", documentationPreset); err != nil {
		return fmt.Errorf("failed to add documentation preset: %w", err)
	}

	// Convert flags to build request
	req := flags.ToBuildRequest()

	// Build the prompt
	result, err := builder.BuildPrompt(req)
	if err != nil {
		return fmt.Errorf("failed to build prompt: %w", err)
	}

	// Format and output the result
	switch flags.OutputFormat {
	case "json":
		// Output as JSON
		jsonData := map[string]any{
			"system_message": result.Prompt.SystemMessage,
			"user_prompt":    result.Prompt.UserPrompt,
			"file_content":   result.Prompt.FileContent,
			"guidelines":     result.Prompt.Guidelines,
		}

		jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		_, err = fmt.Fprintf(output, "%s\n", jsonBytes)
		if err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
	case "text":
		_, err = fmt.Fprintf(output, "%s\n", result.Prompt.String())
		if err != nil {
			return fmt.Errorf("failed to write text output: %w", err)
		}
	default:
		// Default to markdown format
		_, err = fmt.Fprintf(output, "# Generated Prompt\n\n")
		if err != nil {
			return fmt.Errorf("failed to write markdown header: %w", err)
		}

		_, err = fmt.Fprintf(output, "```\n%s\n```\n", result.Prompt.String())
		if err != nil {
			return fmt.Errorf("failed to write markdown content: %w", err)
		}
	}

	return nil
}
