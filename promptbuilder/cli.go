package promptbuilder

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
)

const (
	defaultMaxFileSize = 1024 * 1024 // 1MB default max file size
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
	flagSet.StringVar(&flags.Image, "img", "", "Base64 encoded image data")
	flagSet.StringVar(&flags.Image, "image", "", "Base64 encoded image data")

	// Parse the flags
	err := flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Validate the flags
	err = flags.Validate()
	if err != nil {
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
  -img, --image BASE64      Base64 encoded image data
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
	allowedExtensions := []string{".png"}

	fileProcessor := NewFileProcessor(defaultMaxFileSize, allowedExtensions)

	// Create prompt builder
	builder := New(fileProcessor)

	// Add some default system presets
	codingPreset := "You are an expert software developer. Write clean, efficient, and well-documented code."

	err = builder.AddSystemPreset("coding", codingPreset)
	if err != nil {
		return fmt.Errorf("failed to add coding preset: %w", err)
	}

	analysisPreset := "You are an expert code analyst. Provide detailed analysis and insights."

	err = builder.AddSystemPreset("analysis", analysisPreset)
	if err != nil {
		return fmt.Errorf("failed to add analysis preset: %w", err)
	}

	documentationPreset := "You are an expert technical writer. Create clear and comprehensive documentation."

	err = builder.AddSystemPreset("documentation", documentationPreset)
	if err != nil {
		return fmt.Errorf("failed to add documentation preset: %w", err)
	}

	// Convert flags to build request
	req, err := flags.ToBuildRequest()
	if err != nil {
		return fmt.Errorf("failed to convert flags to build request: %w", err)
	}

	// Build the prompt
	result, err := builder.BuildPrompt(req)
	if err != nil {
		return fmt.Errorf("failed to build prompt: %w", err)
	}

	return formatAndWriteOutput(output, flags.OutputFormat, result.Prompt)
}

// formatAndWriteOutput formats the prompt and writes it to the output writer.
func formatAndWriteOutput(output io.Writer, format string, prompt *Prompt) error {
	var err error // Declare err here

	switch format {
	case "json":
		jsonData := map[string]any{
			"system_message": prompt.SystemMessage,
			"user_prompt":    prompt.UserPrompt,
			"file_content":   prompt.FileContent,
			"guidelines":     prompt.Guidelines,
		}

		jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		_, err = fmt.Fprintf(output, "%s\\n", jsonBytes)
		if err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
	case "text":
		_, err = fmt.Fprintf(output, "%s\\n", prompt.String())
		if err != nil {
			return fmt.Errorf("failed to write text output: %w", err)
		}
	default:
		// Default to markdown format
		_, err = fmt.Fprintf(output, "# Generated Prompt\\n\\n")
		if err != nil {
			return fmt.Errorf("failed to write markdown header: %w", err)
		}

		_, err = fmt.Fprintf(output, "```\\n%s\\n```\\n", prompt.String())
		if err != nil {
			return fmt.Errorf("failed to write markdown content: %w", err)
		}
	}

	return nil
}
