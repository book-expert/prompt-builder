// ./promptbuilder/builder.go
package promptbuilder

import (
	"fmt"
	"strings"
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
		return fmt.Errorf("preset name cannot be empty")
	}

	b.systemPresets[name] = message

	return nil
}

// BuildPrompt constructs a prompt from a BuildRequest.
func (b *Builder) BuildPrompt(req *BuildRequest) (*BuildResult, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build request: %w", err)
	}

	prompt := &Prompt{
		UserPrompt: req.Prompt,
		Guidelines: req.Guidelines,
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
	}

	return &BuildResult{
		Prompt: prompt,
		Error:  nil,
	}, nil
}
