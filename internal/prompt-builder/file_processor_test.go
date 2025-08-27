package promptbuilder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileProcessor_ProcessFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")
	testContent := []byte("func main() {\n    fmt.Println(\"Hello, World!\")\n}")
	
	err := os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	fp := NewFileProcessor(10*1024*1024, []string{".go", ".py", ".js"})
	
	content, err := fp.ProcessFile(testFile)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}
	
	if content == nil {
		t.Fatal("Expected non-nil file content")
	}
	
	if content.Path != testFile {
		t.Errorf("Expected path %s, got %s", testFile, content.Path)
	}
	
	if string(content.Content) != string(testContent) {
		t.Errorf("Expected content %s, got %s", string(testContent), string(content.Content))
	}
	
	if content.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), content.Size)
	}
}

func TestFileProcessor_ProcessFile_NonExistent(t *testing.T) {
	fp := NewFileProcessor(10*1024*1024, []string{".go"})
	
	_, err := fp.ProcessFile("nonexistent.go")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestFileProcessor_ProcessFile_InvalidExtension(t *testing.T) {
	// Create a temporary test file with invalid extension
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("test content")
	
	err := os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	fp := NewFileProcessor(10*1024*1024, []string{".go", ".py"})
	
	_, err = fp.ProcessFile(testFile)
	if err == nil {
		t.Error("Expected error for invalid file extension")
	}
}

func TestFileProcessor_ProcessFile_TooLarge(t *testing.T) {
	// Create a temporary test file that's too large
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "large.go")
	
	// Create content larger than 1KB limit
	largeContent := make([]byte, 2*1024)
	for i := range largeContent {
		largeContent[i] = 'a'
	}
	
	err := os.WriteFile(testFile, largeContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	fp := NewFileProcessor(1024, []string{".go"}) // 1KB limit
	
	_, err = fp.ProcessFile(testFile)
	if err == nil {
		t.Error("Expected error for file too large")
	}
}

func TestFileProcessor_FenceContent(t *testing.T) {
	fp := NewFileProcessor(10*1024*1024, []string{".go"})
	
	content := []byte("func main() {\n    fmt.Println(\"Hello\")\n}")
	fenced := fp.FenceContent(content, "test.go")
	
	if fenced == "" {
		t.Error("Expected non-empty fenced content")
	}
	
	// Should contain BEGIN and END markers
	if !contains(fenced, "BEGIN test.go") {
		t.Error("Expected BEGIN marker in fenced content")
	}
	if !contains(fenced, "END test.go") {
		t.Error("Expected END marker in fenced content")
	}
	
	// Should contain the original content
	if !contains(fenced, "func main()") {
		t.Error("Expected original content in fenced content")
	}
}

func TestFileProcessor_ValidateFile(t *testing.T) {
	fp := NewFileProcessor(10*1024*1024, []string{".go", ".py"})
	
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid go file",
			path:    "test.go",
			wantErr: false,
		},
		{
			name:    "valid py file",
			path:    "script.py",
			wantErr: false,
		},
		{
			name:    "invalid extension",
			path:    "test.txt",
			wantErr: true,
		},
		{
			name:    "no extension",
			path:    "test",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fp.ValidateFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileProcessor.ValidateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
