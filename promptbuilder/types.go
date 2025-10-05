package promptbuilder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Static errors for validation.
var (
	ErrPromptRequired      = errors.New("prompt is required")
	ErrFilePathRequired    = errors.New("file path is required")
	ErrFileContentRequired = errors.New("file content is required")
)

// BuildRequest represents a request to build a prompt. This struct is the main
// data structure that is passed to the prompt builder to construct a prompt.
type BuildRequest struct {
	Prompt        string `json:"prompt"`
	File          string `json:"file,omitempty"`
	Task          string `json:"task,omitempty"`
	SystemMessage string `json:"systemMessage,omitempty"`
	Guidelines    string `json:"guidelines,omitempty"`
	Image         []byte `json:"image,omitempty"`
	OutputFormat  string `json:"outputFormat,omitempty"`
}

// Validate checks if the build request is valid.
func (r *BuildRequest) Validate() error {
	if strings.TrimSpace(r.Prompt) == "" {
		return ErrPromptRequired
	}

	return nil
}

// Prompt represents the assembled prompt. This struct is the output of the prompt
// builder and contains all the components of the prompt.
type Prompt struct {
	SystemMessage string `json:"systemMessage,omitempty"`
	UserPrompt    string `json:"userPrompt"`
	FileContent   string `json:"fileContent,omitempty"`
	Guidelines    string `json:"guidelines,omitempty"`
}

// String returns the formatted prompt as a string.
func (p *Prompt) String() string {
	var parts []string

	if p.SystemMessage != "" {
		parts = append(parts, p.SystemMessage)
	}

	if p.Guidelines != "" {
		parts = append(parts, "Guidelines:", p.Guidelines)
	}

	if p.FileContent != "" {
		parts = append(parts, "File content:", p.FileContent)
	}

	parts = append(parts, p.UserPrompt)

	return strings.Join(parts, "\n\n")
}

// FileContent represents file content with metadata. This struct is used to pass
// file content and metadata between the file processor and the prompt builder.
type FileContent struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
}

// Validate checks if the file content is valid.
func (fc *FileContent) Validate() error {
	if strings.TrimSpace(fc.Path) == "" {
		return ErrFilePathRequired
	}

	if len(fc.Content) == 0 {
		return ErrFileContentRequired
	}

	return nil
}

// SystemPreset represents a predefined system message. This allows for reusable
// system messages that can be referenced by name.
type SystemPreset struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// BuildResult represents the result of building a prompt. This struct is the
// return value of the BuildPrompt function and contains the assembled prompt.
type BuildResult struct {
	Prompt *Prompt `json:"prompt"`
	Error  error   `json:"error,omitempty"`
}

// CLIFlags represents command line interface flags for the prompt builder. This
// struct is used to parse the command line arguments and convert them into a
// BuildRequest.
type CLIFlags struct {
	Prompt        string `json:"prompt"`
	File          string `json:"file,omitempty"`
	Task          string `json:"task,omitempty"`
	SystemMessage string `json:"systemMessage,omitempty"`
	Guidelines    string `json:"guidelines,omitempty"`
	Image         string `json:"image,omitempty"`
	OutputFormat  string `json:"outputFormat,omitempty"`
}

// Validate checks if the CLI flags are valid.
func (f *CLIFlags) Validate() error {
	if strings.TrimSpace(f.Prompt) == "" {
		return ErrPromptRequired
	}

	return nil
}

// ToBuildRequest converts CLI flags to a BuildRequest.
func (f *CLIFlags) ToBuildRequest() (*BuildRequest, error) {
	var imageData []byte

	if f.Image != "" {
		decoded, err := base64.StdEncoding.DecodeString(f.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}

		imageData = decoded
	}

	return &BuildRequest{
		Prompt:        f.Prompt,
		File:          f.File,
		Task:          f.Task,
		SystemMessage: f.SystemMessage,
		Guidelines:    f.Guidelines,
		Image:         imageData,
		OutputFormat:  f.OutputFormat,
	}, nil
}
