package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MaxFileSize is 50MB
	MaxFileSize = 50 * 1024 * 1024
)

// SupportedExtensions contains all supported file extensions
var SupportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
}

// FileHandler handles file operations
type FileHandler struct{}

// ValidationResult represents file validation result
type ValidationResult struct {
	Valid bool
	Error string
}

// New creates a new FileHandler
func New() *FileHandler {
	return &FileHandler{}
}

// ValidateFile validates a single file
func (h *FileHandler) ValidateFile(path string) ValidationResult {
	// Check for path traversal
	if isPathTraversal(path) {
		return ValidationResult{Valid: false, Error: "invalid path (possible path traversal)"}
	}

	// Check if file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ValidationResult{Valid: false, Error: fmt.Sprintf("file not found: %s", path)}
	}
	if err != nil {
		return ValidationResult{Valid: false, Error: err.Error()}
	}

	// Check if it's a file
	if info.IsDir() {
		return ValidationResult{Valid: false, Error: fmt.Sprintf("path is not a file: %s", path)}
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(path))
	if !SupportedExtensions[ext] {
		return ValidationResult{Valid: false, Error: fmt.Sprintf("unsupported file type: %s", ext)}
	}

	// Check size
	if info.Size() > MaxFileSize {
		return ValidationResult{Valid: false, Error: fmt.Sprintf("file too large: %dMB (max: 50MB)", info.Size()/1024/1024)}
	}

	return ValidationResult{Valid: true}
}

// GetFilesFromPath gets all supported files from a path (file or directory)
func (h *FileHandler) GetFilesFromPath(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path not found: %s", path)
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string
	err = filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(walkPath))
			if SupportedExtensions[ext] {
				files = append(files, walkPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// ValidateFiles validates multiple files
func (h *FileHandler) ValidateFiles(paths []string) ([]string, []ValidationResult) {
	var valid, invalid []string
	var invalidResults []ValidationResult

	for _, path := range paths {
		result := h.ValidateFile(path)
		if result.Valid {
			valid = append(valid, path)
		} else {
			invalid = append(invalid, path)
			invalidResults = append(invalidResults, result)
		}
	}

	return valid, invalidResults
}

// ReadURLListFile reads URLs from a file
func (h *FileHandler) ReadURLListFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" && (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return urls, nil
}

// WriteJSONOutput writes data to a JSON file
func (h *FileHandler) WriteJSONOutput(data interface{}, outputPath string) error {
	// Check for path traversal
	if isPathTraversal(outputPath) {
		return fmt.Errorf("invalid output path (possible path traversal): %s", outputPath)
	}

	// Create parent directory
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

// GetMimeType gets the MIME type of a file
func (h *FileHandler) GetMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".bmp":
		return "image/bmp"
	case ".tiff", ".tif":
		return "image/tiff"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	default:
		return "application/octet-stream"
	}
}

// FormatFileSize formats file size for display
func (h *FileHandler) FormatFileSize(sizeBytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(sizeBytes)
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

// IsValidURL checks if a URL is valid
func IsValidURL(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// Check for dangerous schemes
	lowerURL := strings.ToLower(url)
	dangerousSchemes := []string{"javascript:", "file:", "data:", "vbscript:", "mailto:"}
	for _, scheme := range dangerousSchemes {
		if strings.HasPrefix(lowerURL, scheme) {
			return false
		}
	}

	// Check for embedded credentials
	if strings.Contains(url, "@") && strings.Contains(url, "://") {
		schemeEnd := strings.Index(url, "://") + 3
		authorityEnd := strings.Index(url[schemeEnd:], "/")
		if authorityEnd == -1 {
			authorityEnd = len(url)
		} else {
			authorityEnd += schemeEnd
		}
		authority := url[schemeEnd:authorityEnd]
		if strings.Contains(authority, "@") {
			return false
		}
	}

	return true
}

func isPathTraversal(path string) bool {
	if strings.Contains(path, "..") {
		return true
	}
	return false
}
