package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test ---

func TestPromptRunnerRun(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	topic = "test-topic"
	templatePath = filepath.Join(tmpDir, "template.yaml")
	configPath = filepath.Join(tmpDir, "config.yaml")
	outputPath = filepath.Join(tmpDir, "generated.txt")

	// Create dummy template/config files
	writeFile(t, templatePath, "template: mock")
	writeFile(t, configPath, "config: mock")

	// Capture output
	var out bytes.Buffer
	runner := &PromptRunner{Out: &out}

	// Override buildPrompt (optional – here we use the real one)
	// This depends on prompt.Build just writing the file

	// Act
	runner.Run()

	// Assert
	output := out.String()
	assert.Contains(t, output, "Prompt saved to: "+outputPath)

	// Validate prompt file is created
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.NotEmpty(t, data) // Should contain something written by prompt.Build
}

// --- Helpers ---

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
}
