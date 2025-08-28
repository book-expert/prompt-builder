package promptbuilder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileProcessor_validatePathSecurity(t *testing.T) {
	fileProcessor := NewFileProcessor(1024*1024, []string{".go", ".txt"})

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a temporary test file
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

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid file path",
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name:    "path traversal attempt",
			path:    filepath.Join(cwd, "..", "..", "etc", "passwd"),
			wantErr: true,
		},
		{
			name:    "suspicious pattern with ..",
			path:    filepath.Join(cwd, "test", "..", "file.go"),
			wantErr: true,
		},
		{
			name:    "suspicious pattern with /etc",
			path:    cwd + "/test/etc/passwd",
			wantErr: true,
		},
		{
			name:    "outside working directory",
			path:    "/etc/passwd",
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(cwd, "non_existent_file.go"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			absPath, err := filepath.Abs(tt.path)
			if err != nil {
				// If we can't get absolute path, that's also a security issue
				if !tt.wantErr {
					t.Errorf("Expected no error for path %s, but got error: %v", tt.path, err)
				}

				return
			}

			err = fileProcessor.validatePathSecurity(absPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePathSecurity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileProcessor_ProcessFile_Security(t *testing.T) {
	fileProcessor := NewFileProcessor(1024*1024, []string{".go", ".txt"})

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a temporary test file
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
		path    string
		wantErr bool
	}{
		{
			name:    "valid file",
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name:    "path traversal attempt",
			path:    filepath.Join(cwd, "..", "..", "etc", "passwd"),
			wantErr: true,
		},
		{
			name:    "suspicious pattern",
			path:    cwd + "/test/etc/passwd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fileProcessor.ProcessFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
