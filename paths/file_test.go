package paths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFormatListEmpty verifies FormatList returns an empty string for an empty slice.
func TestFormatListEmpty(t *testing.T) {
	result := FormatList([]string{})
	if result != "" {
		t.Errorf("Expected empty string for empty slice, got %q", result)
	}
}

// TestFormatListNonEmpty verifies FormatList produces the expected formatted list.
func TestFormatListNonEmpty(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	expected := "\n- apple\n- banana\n- cherry"
	result := FormatList(items)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestEnsureDirectoryExistsSuccess tests the success branch of EnsureDirectoryExists.
func TestEnsureDirectoryExistsSuccess(t *testing.T) {
	// Save the original mkdirAll function and restore it after the test.
	originalMkdirAll := mkdirAll
	defer func() { mkdirAll = originalMkdirAll }()

	called := false
	var calledPath string
	mkdirAll = func(path string, perm os.FileMode) error {
		called = true
		calledPath = path
		return nil
	}

	testFilePath := filepath.Join("some", "dir", "file.txt")
	EnsureDirectoryExists(testFilePath)

	if !called {
		t.Error("Expected mkdirAll to be called")
	}

	expectedDir := filepath.Dir(testFilePath)
	if calledPath != expectedDir {
		t.Errorf("Expected mkdirAll to be called with %q, got %q", expectedDir, calledPath)
	}
}

// TestEnsureDirectoryExistsError tests the error branch of EnsureDirectoryExists.
func TestEnsureDirectoryExistsError(t *testing.T) {
	// Save original functions and restore them later.
	originalMkdirAll := mkdirAll
	originalFatalf := fatalf
	defer func() {
		mkdirAll = originalMkdirAll
		fatalf = originalFatalf
	}()

	// Override mkdirAll to return an error.
	mkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("fake error")
	}

	// Override fatalf to capture the error message and panic instead of calling os.Exit.
	var fatalMsg string
	fatalf = func(format string, v ...interface{}) {
		fatalMsg = fmt.Sprintf(format, v...)
		panic("fatal")
	}

	testFilePath := filepath.Join("err", "dir", "file.txt")
	// Use recover to catch the panic from fatalf.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected EnsureDirectoryExists to panic")
		} else {
			// Verify that the captured fatal message contains the expected directory.
			expectedSubstr := fmt.Sprintf("Error creating directory %s:", filepath.Dir(testFilePath))
			if !strings.Contains(fatalMsg, expectedSubstr) {
				t.Errorf("Expected fatal message to contain %q, got %q", expectedSubstr, fatalMsg)
			}
		}
	}()

	EnsureDirectoryExists(testFilePath)
	// If no panic occurs, the test will fail in the deferred function.
}
