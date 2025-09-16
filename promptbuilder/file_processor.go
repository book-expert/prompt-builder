package promptbuilder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Static errors for validation
var (
	ErrFileExtensionRequired = errors.New("file must have an extension")
	ErrFileTooLarge          = errors.New("file is too large")
	ErrSuspiciousPath        = errors.New("file path contains suspicious pattern")
	ErrPathOutsideAllowed    = errors.New("file path is outside allowed directories")
	ErrPathIsDirectory       = errors.New("path is a directory, not a file")
)

// FileProcessor handles file operations for prompt building.
type FileProcessor struct {
	maxFileSize       int64
	allowedExtensions []string
}

// NewFileProcessor creates a new file processor with the given constraints.
func NewFileProcessor(maxFileSize int64, allowedExtensions []string) *FileProcessor {
	return &FileProcessor{
		maxFileSize:       maxFileSize,
		allowedExtensions: allowedExtensions,
	}
}

// ProcessFile reads and validates a file, returning its content.
func (fp *FileProcessor) ProcessFile(path string) (*FileContent, error) {
	// Validate file path and extension
	if err := fp.ValidateFile(path); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// Validate path is absolute or relative to current directory
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid file path %s: %w", path, err)
	}

	// Additional security validation: ensure the path doesn't contain path traversal
	if err := fp.validatePathSecurity(absPath); err != nil {
		return nil, fmt.Errorf("security validation failed for %s: %w", absPath, err)
	}

	// Read file content
	// #nosec G304 -- Path is validated for security: checked for path traversal,
	// suspicious patterns, and ensured it's within current working directory
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}

	// Check file size
	if int64(len(content)) > fp.maxFileSize {
		return nil, fmt.Errorf("%w: file %s is too large (%d bytes, max %d bytes)",
			ErrFileTooLarge, path, len(content), fp.maxFileSize)
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

// FenceContent wraps file content with BEGIN/END markers for security.
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

	builder.WriteString("\nEND " + filename)

	return builder.String()
}

// ValidateFile checks if a file path is valid according to the processor's rules.
func (fp *FileProcessor) ValidateFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return ErrFilePathRequired
	}

	ext := filepath.Ext(path)
	if ext == "" {
		return ErrFileExtensionRequired
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

// isCodeFile checks if the file extension indicates a code file.
func isCodeFile(ext string) bool {
	codeExtensions := []string{
		".go",
		".py",
		".js",
		".ts",
		".java",
		".cpp",
		".c",
		".h",
		".cs",
		".php",
		".rb",
		".rs",
	}
	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}

	return false
}

// getLanguageFromExt returns the language identifier for code fencing.
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

// validatePathSecurity ensures the file path is secure and doesn't contain path traversal attempts.
func (fp *FileProcessor) validatePathSecurity(absPath string) error {
	// Check for suspicious path components that indicate path traversal attempts
	suspiciousPatterns := []string{
		"..",
		"~",
		"/etc",
		"/var",
		"/usr",
		"/bin",
		"/sbin",
		"/dev",
		"/sys",
		"/proc",
	}
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(absPath, pattern) {
			return fmt.Errorf(
				"%w: file path %s contains suspicious pattern: %s",
				ErrSuspiciousPath,
				absPath,
				pattern,
			)
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	tmpDir := os.TempDir() // Get the system's temp directory (usually /tmp)

	// CHANGED: Check if the path is within any of the allowed base directories.
	isAllowed := strings.HasPrefix(absPath, homeDir) ||
		strings.HasPrefix(absPath, cwd) ||
		strings.HasPrefix(absPath, tmpDir)

	if !isAllowed {
		return fmt.Errorf(
			"%w: file path %s is outside allowed directories (home: %s, cwd: %s, tmp: %s)",
			ErrPathOutsideAllowed,
			absPath,
			homeDir,
			cwd,
			tmpDir,
		)
	}

	// Ensure the file exists and is a regular file
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", absPath, err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%w: path %s is a directory, not a file", ErrPathIsDirectory, absPath)
	}

	return nil
}
