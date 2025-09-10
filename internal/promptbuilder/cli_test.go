package promptbuilder

import (
	"bytes"
	"os"
	"testing"
)

func TestCLIFlags_Validate(t *testing.T) {
	tests := []struct {
		name    string
		flags   CLIFlags
		wantErr bool
	}{
		{
			name: "valid flags with prompt only",
			flags: CLIFlags{
				Prompt: "test prompt",
			},
			wantErr: false,
		},
		{
			name: "valid flags with prompt and file",
			flags: CLIFlags{
				Prompt: "test prompt",
				File:   "test.go",
			},
			wantErr: false,
		},
		{
			name: "empty prompt should fail",
			flags: CLIFlags{
				Prompt: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace-only prompt should fail",
			flags: CLIFlags{
				Prompt: "   \n\t  ",
			},
			wantErr: true,
		},
		{
			name: "valid flags with task preset",
			flags: CLIFlags{
				Prompt: "test prompt",
				Task:   "coding",
			},
			wantErr: false,
		},
		{
			name: "valid flags with system message",
			flags: CLIFlags{
				Prompt:        "test prompt",
				SystemMessage: "You are a helpful assistant",
			},
			wantErr: false,
		},
		{
			name: "valid flags with guidelines",
			flags: CLIFlags{
				Prompt:     "test prompt",
				Guidelines: "Follow best practices",
			},
			wantErr: false,
		},

		{
			name: "valid flags with output format",
			flags: CLIFlags{
				Prompt:       "test prompt",
				OutputFormat: "json",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.flags.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CLIFlags.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIFlags_ToBuildRequest(t *testing.T) {
	flags := CLIFlags{
		Prompt:        "test prompt",
		File:          "test.go",
		Task:          "coding",
		SystemMessage: "You are a helpful assistant",
		Guidelines:    "Follow best practices",
		OutputFormat:  "json",
	}

	req := flags.ToBuildRequest()

	if req.Prompt != flags.Prompt {
		t.Errorf("Expected prompt %s, got %s", flags.Prompt, req.Prompt)
	}

	if req.File != flags.File {
		t.Errorf("Expected file %s, got %s", flags.File, req.File)
	}

	if req.Task != flags.Task {
		t.Errorf("Expected task %s, got %s", flags.Task, req.Task)
	}

	if req.SystemMessage != flags.SystemMessage {
		t.Errorf("Expected system message %s, got %s", flags.SystemMessage, req.SystemMessage)
	}

	if req.Guidelines != flags.Guidelines {
		t.Errorf("Expected guidelines %s, got %s", flags.Guidelines, req.Guidelines)
	}

	if req.OutputFormat != flags.OutputFormat {
		t.Errorf("Expected output format %s, got %s", flags.OutputFormat, req.OutputFormat)
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    CLIFlags
		wantErr bool
	}{
		{
			name: "basic prompt",
			args: []string{"-p", "test prompt"},
			want: CLIFlags{
				Prompt: "test prompt",
			},
			wantErr: false,
		},
		{
			name: "prompt with file",
			args: []string{"-p", "test prompt", "-f", "test.go"},
			want: CLIFlags{
				Prompt: "test prompt",
				File:   "test.go",
			},
			wantErr: false,
		},
		{
			name: "prompt with task",
			args: []string{"-p", "test prompt", "-t", "coding"},
			want: CLIFlags{
				Prompt: "test prompt",
				Task:   "coding",
			},
			wantErr: false,
		},
		{
			name: "prompt with system message",
			args: []string{"-p", "test prompt", "-sys", "You are a helpful assistant"},
			want: CLIFlags{
				Prompt:        "test prompt",
				SystemMessage: "You are a helpful assistant",
			},
			wantErr: false,
		},
		{
			name: "prompt with guidelines",
			args: []string{"-p", "test prompt", "-g", "Follow best practices"},
			want: CLIFlags{
				Prompt:     "test prompt",
				Guidelines: "Follow best practices",
			},
			wantErr: false,
		},
		{
			name: "prompt with output format",
			args: []string{"-p", "test prompt", "-o", "json"},
			want: CLIFlags{
				Prompt:       "test prompt",
				OutputFormat: "json",
			},
			wantErr: false,
		},
		{
			name:    "missing prompt should fail",
			args:    []string{"-f", "test.go"},
			wantErr: true,
		},
		{
			name:    "empty prompt should fail",
			args:    []string{"-p", ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, err := ParseFlags(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if flags.Prompt != tt.want.Prompt {
					t.Errorf("Expected prompt %s, got %s", tt.want.Prompt, flags.Prompt)
				}

				if flags.File != tt.want.File {
					t.Errorf("Expected file %s, got %s", tt.want.File, flags.File)
				}

				if flags.Task != tt.want.Task {
					t.Errorf("Expected task %s, got %s", tt.want.Task, flags.Task)
				}

				if flags.SystemMessage != tt.want.SystemMessage {
					t.Errorf(
						"Expected system message %s, got %s",
						tt.want.SystemMessage,
						flags.SystemMessage,
					)
				}

				if flags.Guidelines != tt.want.Guidelines {
					t.Errorf("Expected guidelines %s, got %s", tt.want.Guidelines, flags.Guidelines)
				}

				if flags.OutputFormat != tt.want.OutputFormat {
					t.Errorf(
						"Expected output format %s, got %s",
						tt.want.OutputFormat,
						flags.OutputFormat,
					)
				}
			}
		})
	}
}

func TestRunCLI(t *testing.T) {
	// Create a temporary test file in current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tmpFile, err := os.CreateTemp(cwd, "test_*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	// Write some test content
	testContent := `package main

func main() {
    fmt.Println("Hello, World!")
}`
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "basic prompt",
			args:    []string{"-p", "Explain this code"},
			wantErr: false,
		},
		{
			name:    "prompt with file",
			args:    []string{"-p", "Explain this code", "-f", tmpFile.Name()},
			wantErr: false,
		},
		{
			name:    "prompt with task",
			args:    []string{"-p", "Explain this code", "-t", "coding"},
			wantErr: false,
		},
		{
			name:    "missing prompt should fail",
			args:    []string{"-f", tmpFile.Name()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RunCLI(tt.args, &buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCLI() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && buf.Len() == 0 {
				t.Error("Expected output, got empty buffer")
			}
		})
	}
}
