package promptbuilder

import (
	"testing"
)

func TestBuildRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     BuildRequest
		wantErr bool
	}{
		{
			name: "valid request with prompt",
			req: BuildRequest{
				Prompt: "test prompt",
			},
			wantErr: false,
		},
		{
			name: "empty prompt should fail",
			req: BuildRequest{
				Prompt: "",
			},
			wantErr: true,
		},
		{
			name: "valid request with file",
			req: BuildRequest{
				Prompt: "test prompt",
				File:   "test.go",
			},
			wantErr: false,
		},
		{
			name: "valid request with task",
			req: BuildRequest{
				Prompt: "test prompt",
				Task:   "coding",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrompt_String(t *testing.T) {
	prompt := Prompt{
		SystemMessage: "You are a helpful assistant",
		UserPrompt:    "Explain this code",
		FileContent:   "func main() { }",
	}

	result := prompt.String()

	if result == "" {
		t.Error("Expected non-empty prompt string")
	}

	// Should contain all components
	if !contains(result, "You are a helpful assistant") {
		t.Error("Expected system message in prompt")
	}

	if !contains(result, "Explain this code") {
		t.Error("Expected user prompt in prompt")
	}

	if !contains(result, "func main() { }") {
		t.Error("Expected file content in prompt")
	}
}

func TestFileContent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		content FileContent
		wantErr bool
	}{
		{
			name: "valid file content",
			content: FileContent{
				Path:    "test.go",
				Content: []byte("func main() { }"),
			},
			wantErr: false,
		},
		{
			name: "empty path should fail",
			content: FileContent{
				Path:    "",
				Content: []byte("content"),
			},
			wantErr: true,
		},
		{
			name: "empty content should fail",
			content: FileContent{
				Path:    "test.go",
				Content: []byte{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.content.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileContent.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}

				return false
			}()))
}
