package promptbuilder

import (
	"fmt"
	"strings"
)

// BuildRequest represents a request to build a prompt
type BuildRequest struct {
	Prompt        string `json:"prompt"`
	File          string `json:"file,omitempty"`
	Task          string `json:"task,omitempty"`
	SystemMessage string `json:"system_message,omitempty"`
	Guidelines    string `json:"guidelines,omitempty"`
	Image         string `json:"image,omitempty"`
	OutputFormat  string `json:"output_format,omitempty"`
	EstimateTokens bool  `json:"estimate_tokens,omitempty"`
	Model         string `json:"model,omitempty"`
}

// Validate checks if the build request is valid
func (r *BuildRequest) Validate() error {
	if strings.TrimSpace(r.Prompt) == "" {
		return fmt.Errorf("prompt is required")
	}
	return nil
}

// Prompt represents the assembled prompt
type Prompt struct {
	SystemMessage string `json:"system_message,omitempty"`
	UserPrompt    string `json:"user_prompt"`
	FileContent   string `json:"file_content,omitempty"`
	Guidelines    string `json:"guidelines,omitempty"`
	ImagePath     string `json:"image_path,omitempty"`
	TokenEstimate int    `json:"token_estimate,omitempty"`
}

// String returns the formatted prompt as a string
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

// FileContent represents file content with metadata
type FileContent struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
}

// Validate checks if the file content is valid
func (fc *FileContent) Validate() error {
	if strings.TrimSpace(fc.Path) == "" {
		return fmt.Errorf("file path is required")
	}
	if len(fc.Content) == 0 {
		return fmt.Errorf("file content is required")
	}
	return nil
}

// SystemPreset represents a predefined system message
type SystemPreset struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// BuildResult represents the result of building a prompt
type BuildResult struct {
	Prompt        *Prompt `json:"prompt"`
	TokenEstimate int     `json:"token_estimate,omitempty"`
	Error         error   `json:"error,omitempty"`
}
