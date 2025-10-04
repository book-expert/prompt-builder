package promptbuilder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/book-expert/prompt-builder/promptbuilder"
)

// setupFileProcessorTest creates a temporary file and returns its path,
// along with the current working directory and a cleanup function.
func setupFileProcessorTest(t *testing.T) (string, string, func()) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tmpFile, err := os.CreateTemp(cwd, "test_*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write some test content
	testContent := `package main
func main() {
    fmt.Println("Hello, World!")
}`

	_, err = tmpFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cleanup := func() {
		t.Helper()

		err := os.Remove(tmpFile.Name())
		if err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}

	return tmpFile.Name(), cwd, cleanup
}

func TestFileProcessor_ProcessFile_Security(t *testing.T) {
	t.Parallel()

	fileProcessor := promptbuilder.NewFileProcessor(1024*1024, []string{".go", ".txt"})

	tmpFileName, cwd, cleanup := setupFileProcessorTest(t)
	t.Cleanup(cleanup)

	var processFileSecurityTests = []struct { // Extracted test cases
		name    string
		path    string
		wantErr bool
	}{
		// These paths will be dynamically set in TestFileProcessor_ProcessFile_Security
		// to include cwd and tmpFile.Name()
		{
			name:    "path traversal attempt",
			path:    filepath.Join("..", "..", "etc", "passwd"), // Relative path for dynamic joining
			wantErr: true,
		},
		{
			name:    "suspicious pattern",
			path:    "/test/etc/passwd", // Relative path for dynamic joining
			wantErr: true,
		},
	}

	// Dynamically add test cases that depend on cwd and tmpFile
	dynamicTests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid file",
			path:    tmpFileName,
			wantErr: false,
		},
		{
			name:    "path traversal attempt (dynamic)",
			path:    filepath.Join(cwd, processFileSecurityTests[0].path),
			wantErr: true,
		},
		{
			name:    "suspicious pattern (dynamic)",
			path:    filepath.Join(cwd, processFileSecurityTests[1].path),
			wantErr: true,
		},
	}

	// Combine static and dynamic tests and iterate immediately
	for _, testCase := range append(processFileSecurityTests, dynamicTests...) {
		// Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := fileProcessor.ProcessFile(testCase.path)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ProcessFile() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}
