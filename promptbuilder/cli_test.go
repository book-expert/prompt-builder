package promptbuilder_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/book-expert/prompt-builder/promptbuilder"
)

const (
	// sampleImageB64 is a very small 1x1 PNG used in tests.
	sampleImageB64Part1 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAA"
	sampleImageB64Part2 = "AAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="
)

type validateCase struct {
	name    string
	flags   promptbuilder.CLIFlags
	wantErr bool
}

func runFlagValidationSubtests(t *testing.T, cases []validateCase) {
	t.Helper()

	for _, item := range cases {
		testCase := item
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			validateErr := testCase.flags.Validate()
			if (validateErr != nil) != testCase.wantErr {
				t.Errorf("CLIFlags.Validate() error = %v, wantErr %v", validateErr, testCase.wantErr)
			}
		})
	}
}

// TestCLIFlagsValidateBasic tests basic validation scenarios.
func TestCLIFlagsValidateBasic(t *testing.T) {
	t.Parallel()

	cases := []validateCase{
		{
			name: "valid flags with prompt only",
			flags: promptbuilder.CLIFlags{
				Prompt:        "test prompt",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: false,
		},
		{
			name: "valid flags with prompt and file",
			flags: promptbuilder.CLIFlags{
				Prompt:        "test prompt",
				File:          "test.go",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: false,
		},
	}
	runFlagValidationSubtests(t, cases)
}

// TestCLIFlagsValidateEmptyPrompt tests validation of empty prompts.
func TestCLIFlagsValidateEmptyPrompt(t *testing.T) {
	t.Parallel()

	cases := []validateCase{
		{
			name: "empty prompt should fail",
			flags: promptbuilder.CLIFlags{
				Prompt:        "",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: true,
		},
		{
			name: "whitespace-only prompt should fail",
			flags: promptbuilder.CLIFlags{
				Prompt:        "   \n\t  ",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: true,
		},
	}
	runFlagValidationSubtests(t, cases)
}

// TestCLIFlagsValidateWithPresets tests validation with presets and additional fields.
func TestCLIFlagsValidateWithPresets(t *testing.T) {
	t.Parallel()

	cases := []validateCase{
		{
			name: "valid flags with task preset",
			flags: promptbuilder.CLIFlags{
				Prompt:        "test prompt",
				File:          "",
				Task:          "coding",
				SystemMessage: "",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: false,
		},
		{
			name: "valid flags with system message",
			flags: promptbuilder.CLIFlags{
				Prompt:        "test prompt",
				File:          "",
				Task:          "",
				SystemMessage: "You are a helpful assistant",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: false,
		},
		{
			name: "valid flags with guidelines",
			flags: promptbuilder.CLIFlags{
				Prompt:        "test prompt",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "Follow best practices",
				Image:         "",
				OutputFormat:  "",
			},
			wantErr: false,
		},
		{
			name: "valid flags with output format",
			flags: promptbuilder.CLIFlags{
				Prompt:        "test prompt",
				File:          "",
				Task:          "",
				SystemMessage: "",
				Guidelines:    "",
				Image:         "",
				OutputFormat:  "json",
			},
			wantErr: false,
		},
	}
	runFlagValidationSubtests(t, cases)
}

func TestCLIFlags_ToBuildRequest(t *testing.T) {
	t.Parallel()

	flags := promptbuilder.CLIFlags{
		Prompt:        "test prompt",
		File:          "test.go",
		Task:          "coding",
		SystemMessage: "You are a helpful assistant",
		Guidelines:    "Follow best practices",
		OutputFormat:  "json",
		Image:         sampleImageB64Part1 + sampleImageB64Part2,
	}

	req, err := flags.ToBuildRequest()
	if err != nil {
		t.Fatalf("Failed to convert flags to build request: %v", err)
	}

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

	// Check image data
	decodedImage, err := base64.StdEncoding.DecodeString(flags.Image)
	if err != nil {
		t.Fatalf("Failed to decode image: %v", err)
	}

	if !bytes.Equal(req.Image, decodedImage) {
		t.Errorf("Expected image data %v, got %v", decodedImage, req.Image)
	}
}

func TestParseFlags_Basic(t *testing.T) {
	t.Parallel()

	flags, parseErr := promptbuilder.ParseFlags([]string{"-p", "test prompt"})
	if parseErr != nil {
		t.Fatalf("ParseFlags() unexpected error = %v", parseErr)
	}

	if flags.Prompt != "test prompt" {
		t.Errorf("Expected prompt %s, got %s", "test prompt", flags.Prompt)
	}
}

func TestParseFlags_WithFile(t *testing.T) {
	t.Parallel()

	flags, parseErr := promptbuilder.ParseFlags([]string{"-p", "test prompt", "-f", "test.go"})
	if parseErr != nil {
		t.Fatalf("ParseFlags() unexpected error = %v", parseErr)
	}

	if flags.File != "test.go" {
		t.Errorf("Expected file %s, got %s", "test.go", flags.File)
	}
}

func TestParseFlags_WithTask(t *testing.T) {
	t.Parallel()

	flags, parseErr := promptbuilder.ParseFlags([]string{"-p", "test prompt", "-t", "coding"})
	if parseErr != nil {
		t.Fatalf("ParseFlags() unexpected error = %v", parseErr)
	}

	if flags.Task != "coding" {
		t.Errorf("Expected task %s, got %s", "coding", flags.Task)
	}
}

func TestParseFlags_WithSystemMessage(t *testing.T) {
	t.Parallel()

	sys := "You are a helpful assistant"

	flags, parseErr := promptbuilder.ParseFlags([]string{"-p", "test prompt", "-sys", sys})
	if parseErr != nil {
		t.Fatalf("ParseFlags() unexpected error = %v", parseErr)
	}

	if flags.SystemMessage != sys {
		t.Errorf("Expected system message %s, got %s", sys, flags.SystemMessage)
	}
}

func TestParseFlags_WithGuidelines(t *testing.T) {
	t.Parallel()

	guide := "Follow best practices"

	flags, parseErr := promptbuilder.ParseFlags([]string{"-p", "test prompt", "-g", guide})
	if parseErr != nil {
		t.Fatalf("ParseFlags() unexpected error = %v", parseErr)
	}

	if flags.Guidelines != guide {
		t.Errorf("Expected guidelines %s, got %s", guide, flags.Guidelines)
	}
}

func TestParseFlags_WithOutputFormat(t *testing.T) {
	t.Parallel()

	flags, parseErr := promptbuilder.ParseFlags([]string{"-p", "test prompt", "-o", "json"})
	if parseErr != nil {
		t.Fatalf("ParseFlags() unexpected error = %v", parseErr)
	}

	if flags.OutputFormat != "json" {
		t.Errorf("Expected output format %s, got %s", "json", flags.OutputFormat)
	}
}

func TestParseFlags_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
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

	for _, item := range tests {
		testCase := item
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := promptbuilder.ParseFlags(testCase.args)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}

func TestRunCLI(t *testing.T) {
	t.Parallel()

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
			name:    "prompt with image",
			args:    []string{"-p", "Explain this image", "-img", sampleImageB64Part1 + sampleImageB64Part2},
			wantErr: false,
		},
		{
			name:    "prompt with task",
			args:    []string{"-p", "Explain this code", "-t", "coding"},
			wantErr: false,
		},
		{
			name:    "missing prompt should fail",
			args:    []string{"-img", sampleImageB64Part1 + sampleImageB64Part2},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			err := promptbuilder.RunCLI(testCase.args, &buf)
			if (err != nil) != testCase.wantErr {
				t.Errorf("RunCLI() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !testCase.wantErr && buf.Len() == 0 {
				t.Error("Expected output, got empty buffer")
			}
		})
	}
}
