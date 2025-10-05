package promptbuilder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrFileExtensionRequired is returned when a file path doesn't have an extension.
var (
	ErrFileExtensionRequired   = errors.New("file must have an extension")
	ErrFileTooLarge            = errors.New("file is too large")
	ErrSuspiciousPath          = errors.New("file path contains suspicious pattern")
	ErrPathOutsideAllowed      = errors.New("file path is outside allowed directories")
	ErrPathIsDirectory         = errors.New("path is a directory, not a file")
	ErrFileExtensionNotAllowed = errors.New("file extension is not allowed") // Add this line
)

// FileProcessor handles file operations for prompt building. It is responsible for
// reading, validating, and fencing file content to be included in a prompt.
type FileProcessor struct {
	maxFileSize       int64
	allowedExtensions []string
}

// NewFileProcessor creates a new file processor with the given constraints. This
// function is the designated constructor for the FileProcessor struct and ensures
// that the processor is initialized with the necessary constraints.
func NewFileProcessor(maxFileSize int64, allowedExtensions []string) *FileProcessor {
	return &FileProcessor{
		maxFileSize:       maxFileSize,
		allowedExtensions: allowedExtensions,
	}
}

// ProcessFile reads and validates a file, returning its content. This is the main
// entry point for the file processor and is responsible for orchestrating the
// entire file processing workflow.
func (fp *FileProcessor) ProcessFile(path string) (*FileContent, error) {
	// Validate file path and extension
	err := fp.ValidateFile(path)
	if err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// Validate path is absolute or relative to current directory
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid file path %s: %w", path, err)
	}

	// Additional security validation: ensure the path doesn't contain path traversal
	err = fp.validatePathSecurity(absPath)
	if err != nil {
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

// FenceContent wraps file content with BEGIN/END markers for security and clarity.
// This makes it clear to the model where the file content begins and ends.
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
// This function is responsible for ensuring that the file path is not empty, has a
// valid extension, and that the extension is allowed.
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
		return fmt.Errorf("%w: file extension %s is not allowed. Allowed extensions: %v",
			ErrFileExtensionNotAllowed, // Use the static error
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

// validatePathSecurity ensures the file path is secure and doesn't contain path
// traversal attempts. This function is a critical security measure to prevent
// the model from accessing unauthorized files.
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
