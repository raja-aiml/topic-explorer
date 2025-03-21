package paths

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Declare package-level variables for dependency injection.
var mkdirAll = os.MkdirAll
var fatalf = log.Fatalf

// EnsureDirectoryExists checks if a directory exists and creates it if necessary.
func EnsureDirectoryExists(filePath string) {
	dir := filepath.Dir(filePath)
	if err := mkdirAll(dir, os.ModePerm); err != nil {
		fatalf("Error creating directory %s: %v", dir, err)
	}
}

// FormatList converts a slice of strings into a formatted list.
func FormatList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return "\n- " + strings.Join(items, "\n- ")
}
