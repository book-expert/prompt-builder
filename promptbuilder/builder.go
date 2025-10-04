// Package promptbuilder provides functionality for building prompts from various components
// including files, system messages, and guidelines.
package promptbuilder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// ErrPresetNameEmpty is returned when trying to add a system preset with an empty name.
var (
	ErrPresetNameEmpty = errors.New("preset name cannot be empty")
)

// Builder is the main engine for constructing prompts.
type Builder struct {
	fileProcessor *FileProcessor
	systemPresets map[string]string
}

// New creates a new prompt builder with a given file processor.
func New(fp *FileProcessor) *Builder {
	return &Builder{
		fileProcessor: fp,
		systemPresets: make(map[string]string),
	}
}

// AddSystemPreset adds a named system message preset to the builder.
func (b *Builder) AddSystemPreset(name, message string) error {
	if strings.TrimSpace(name) == "" {
		return ErrPresetNameEmpty // Use the static error
	}

	b.systemPresets[name] = message

	return nil
}

// BuildPrompt constructs a prompt from a BuildRequest.
func (b *Builder) BuildPrompt(req *BuildRequest) (*BuildResult, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid build request: %w", err)
	}

	prompt := &Prompt{
		UserPrompt:    req.Prompt,
		Guidelines:    req.Guidelines,
		SystemMessage: "", // Initialize SystemMessage
		FileContent:   "", // Initialize FileContent
	}

	// Handle the system message logic
	if req.SystemMessage != "" {
		prompt.SystemMessage = req.SystemMessage
	} else if req.Task != "" {
		if preset, ok := b.systemPresets[req.Task]; ok {
			prompt.SystemMessage = preset
		}
	}

	// Handle the file content
	if req.File != "" {
		fileContent, err := b.fileProcessor.ProcessFile(req.File)
		if err != nil {
			return nil, fmt.Errorf("failed to process file: %w", err)
		}

		prompt.FileContent = b.fileProcessor.FenceContent(fileContent.Content, fileContent.Path)
	} else if len(req.Image) > 0 {
		// Assuming image is PNG for now, as per png-to-text-service context
		encodedImage := base64.StdEncoding.EncodeToString(req.Image)
		prompt.FileContent = b.fileProcessor.FenceContent([]byte("data:image/png;base64,"+encodedImage), "image.png")
	}

	return &BuildResult{
		Prompt: prompt,
		Error:  nil,
	}, nil
}
