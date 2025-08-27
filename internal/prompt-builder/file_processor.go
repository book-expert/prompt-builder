package promptbuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileProcessor handles file operations for prompt building
type FileProcessor struct {
	maxFileSize       int64
	allowedExtensions []string
}

// NewFileProcessor creates a new file processor with the given constraints
func NewFileProcessor(maxFileSize int64, allowedExtensions []string) *FileProcessor {
	return &FileProcessor{
		maxFileSize:       maxFileSize,
		allowedExtensions: allowedExtensions,
	}
}

// ProcessFile reads and validates a file, returning its content
func (fp *FileProcessor) ProcessFile(path string) (*FileContent, error) {
	// Validate file path and extension
	if err := fp.ValidateFile(path); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}
	
	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	
	// Check file size
	if int64(len(content)) > fp.maxFileSize {
		return nil, fmt.Errorf("file %s is too large (%d bytes, max %d bytes)", 
			path, len(content), fp.maxFileSize)
	}
	
	// Get file info for size
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}
	
	return &FileContent{
		Path:    path,
		Content: content,
		Size:    fileInfo.Size(),
	}, nil
}

// FenceContent wraps file content with BEGIN/END markers for security
func (fp *FileProcessor) FenceContent(content []byte, filename string) string {
	ext := filepath.Ext(filename)
	
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("BEGIN %s\n", filename))
	
	// Add code fence if it's a code file
	if isCodeFile(ext) {
		builder.WriteString(fmt.Sprintf("```%s\n", getLanguageFromExt(ext)))
	}
	
	builder.Write(content)
	
	if isCodeFile(ext) {
		builder.WriteString("\n```")
	}
	
	builder.WriteString(fmt.Sprintf("\nEND %s", filename))
	
	return builder.String()
}

// ValidateFile checks if a file path is valid according to the processor's rules
func (fp *FileProcessor) ValidateFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("file path is required")
	}
	
	ext := filepath.Ext(path)
	if ext == "" {
		return fmt.Errorf("file must have an extension")
	}
	
	// Check if extension is allowed
	allowed := false
	for _, allowedExt := range fp.allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	
	if !allowed {
		return fmt.Errorf("file extension %s is not allowed. Allowed extensions: %v", 
			ext, fp.allowedExtensions)
	}
	
	return nil
}

// isCodeFile checks if the file extension indicates a code file
func isCodeFile(ext string) bool {
	codeExtensions := []string{".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h", ".cs", ".php", ".rb", ".rs"}
	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}
	return false
}

// getLanguageFromExt returns the language identifier for code fencing
func getLanguageFromExt(ext string) string {
	languageMap := map[string]string{
		".go":   "go",
		".py":   "python",
		".js":   "javascript",
		".ts":   "typescript",
		".java": "java",
		".cpp":  "cpp",
		".c":    "c",
		".h":    "c",
		".cs":   "csharp",
		".php":  "php",
		".rb":   "ruby",
		".rs":   "rust",
	}
	
	if lang, exists := languageMap[ext]; exists {
		return lang
	}
	
	return "text"
}
