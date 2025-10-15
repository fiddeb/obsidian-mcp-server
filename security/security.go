package security

import (
"fmt"
"strings"
)

// ValidatePath validates file paths to prevent directory traversal and other attacks
func ValidatePath(path string) error {
	// Check for directory traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: directory traversal not allowed")
	}

	// Check for absolute paths
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid path: absolute paths not allowed")
	}

	// Check for null bytes and other control characters
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("invalid path: contains null byte")
	}

	// Check for common dangerous patterns that could cause issues
	// Note: We allow & and other characters that are valid in Obsidian filenames
	dangerous := []string{"~", "$", "`", "|", ";"}
for _, d := range dangerous {
		if strings.Contains(path, d) {
			return fmt.Errorf("invalid path: contains dangerous character '%s'", d)
}
}

// Validate that path ends with .md
if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, "/") {
// Allow paths without extension only if they're folder references
if strings.Contains(path, ".") {
return fmt.Errorf("invalid path: only .md files are supported")
}
}

return nil
}

// SanitizeContent sanitizes content to prevent injection attacks
func SanitizeContent(content string) string {
// Remove null bytes
content = strings.ReplaceAll(content, "\x00", "")

return content
}
