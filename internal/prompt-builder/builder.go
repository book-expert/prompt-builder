package promptbuilder

import (
	"fmt"
	"strings"
)

// PromptBuilder assembles prompts from various components
type PromptBuilder struct {
	fileProcessor *FileProcessor
	systemPresets map[string]string
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder(fileProcessor *FileProcessor) *PromptBuilder {
	return &PromptBuilder{
		fileProcessor: fileProcessor,
		systemPresets: make(map[string]string),
	}
}

// BuildPrompt assembles a prompt from the given request
func (pb *PromptBuilder) BuildPrompt(req *BuildRequest) (*BuildResult, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Create the prompt
	prompt := &Prompt{
		UserPrompt: req.Prompt,
	}

	// Add system message from task preset if specified
	if req.Task != "" {
		if preset, exists := pb.systemPresets[req.Task]; exists {
			prompt.SystemMessage = preset
		}
	}

	// Override with custom system message if provided
	if req.SystemMessage != "" {
		prompt.SystemMessage = req.SystemMessage
	}

	// Add guidelines if provided
	if req.Guidelines != "" {
		prompt.Guidelines = req.Guidelines
	}

	// Process file if provided
	if req.File != "" {
		fileContent, err := pb.fileProcessor.ProcessFile(req.File)
		if err != nil {
			return nil, fmt.Errorf("failed to process file: %w", err)
		}

		// Fence the content for security
		prompt.FileContent = pb.fileProcessor.FenceContent(fileContent.Content, fileContent.Path)
	}

	// Add image path if provided
	if req.Image != "" {
		prompt.ImagePath = req.Image
	}

	return &BuildResult{
		Prompt: prompt,
	}, nil
}

// AddSystemPreset adds a system message preset for a specific task
func (pb *PromptBuilder) AddSystemPreset(task, message string) error {
	if strings.TrimSpace(task) == "" {
		return fmt.Errorf("task name is required")
	}
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("system message is required")
	}

	pb.systemPresets[task] = message
	return nil
}

// GetSystemPreset retrieves a system message preset for a specific task
func (pb *PromptBuilder) GetSystemPreset(task string) (string, bool) {
	preset, exists := pb.systemPresets[task]
	return preset, exists
}

// ListSystemPresets returns all available system presets
func (pb *PromptBuilder) ListSystemPresets() []string {
	presets := make([]string, 0, len(pb.systemPresets))
	for task := range pb.systemPresets {
		presets = append(presets, task)
	}
	return presets
}
