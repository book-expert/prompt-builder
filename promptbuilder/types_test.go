package promptbuilder_test

import (
	"testing"

	"github.com/book-expert/prompt-builder/promptbuilder"
)

// buildRequestValidationTests returns common test cases for BuildRequest validation.
func buildRequestValidationTests() []struct {
	name    string
	req     promptbuilder.BuildRequest
	wantErr bool
} {
	return []struct {
		name    string
		req     promptbuilder.BuildRequest
		wantErr bool
	}{
		{
			name: "valid request with prompt",
			req: promptbuilder.BuildRequest{
				Prompt:        "test prompt",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         nil,
				OutputFormat:  "",
			},
			wantErr: false,
		},
		{
			name: "empty prompt should fail",
			req: promptbuilder.BuildRequest{
				Prompt:        "",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         nil,
				OutputFormat:  "",
			},
			wantErr: true,
		},
	}
}

// TestBuildRequestValidateBasic tests basic validation scenarios.
func TestBuildRequestValidateBasic(t *testing.T) {
	t.Parallel()

	tests := buildRequestValidationTests()

	for _, testCase := range tests {
		// Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := testCase.req.Validate()
			if (err != nil) != testCase.wantErr {
				t.Errorf("BuildRequest.Validate() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}

// TestBuildRequestValidateWithComponents tests validation with various components.
func TestBuildRequestValidateWithComponents(t *testing.T) {
	t.Parallel()

	tests := []struct { // Extracted test cases
		name    string
		req     promptbuilder.BuildRequest
		wantErr bool
	}{
		{
			name: "valid request with prompt and file",
			req: promptbuilder.BuildRequest{
				Prompt:        "test prompt",
				File:          "test.go",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         nil,
				OutputFormat:  "",
			},
			wantErr: false,
		},
		{
			name: "valid request with task",
			req: promptbuilder.BuildRequest{
				Prompt:        "test prompt",
				File:          "",
				Task:          "coding",
				SystemMessage: "",
				Guidelines:    "",
				Image:         nil,
				OutputFormat:  "",
			},
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		// Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := testCase.req.Validate()
			if (err != nil) != testCase.wantErr {
				t.Errorf("BuildRequest.Validate() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}
// TestPromptString tests the String method of the Prompt struct.
func TestPromptString(t *testing.T) {
	t.Parallel()

	prompt := promptbuilder.Prompt{
		SystemMessage: "System message.",
		UserPrompt:    "User prompt.",
		FileContent:   "File content.",
		Guidelines:    "Guidelines.",
	}

	result := prompt.String()

	if result == "" {
		t.Error("Expected non-empty prompt string")
	}

	// Should contain all components
	if !contains(result, "System message.") {
		t.Error("Expected system message in prompt")
	}

	if !contains(result, "User prompt.") {
		t.Error("Expected user prompt in prompt")
	}

	if !contains(result, "File content.") {
		t.Error("Expected file content in prompt")
	}
}

func TestFileContentValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content promptbuilder.FileContent
		wantErr bool
	}{
		{
			name: "valid file content",
			content: promptbuilder.FileContent{
				Path:    "test.txt",
				Content: []byte("hello"),
				Size:    5,
			},
			wantErr: false,
		},
		{
			name: "empty path should fail",
			content: promptbuilder.FileContent{
				Path:    "",
				Content: []byte("hello"),
				Size:    0,
			},
			wantErr: true,
		},
		{
			name: "empty content should fail",
			content: promptbuilder.FileContent{
				Path:    "test.txt",
				Content: []byte{},
				Size:    0,
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		// Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := testCase.content.Validate()
			if (err != nil) != testCase.wantErr {
				t.Errorf("FileContent.Validate() error = %v, wantErr %v", err, testCase.wantErr)
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