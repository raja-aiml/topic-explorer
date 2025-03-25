package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptRunnerRun(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	topic = "test-topic"
	templatePath = filepath.Join(tmpDir, "template.yaml")
	configPath = filepath.Join(tmpDir, "config.yaml")
	outputPath = filepath.Join(tmpDir, "generated.txt")

	// Create a dummy template file with valid YAML structure.
	// This file will be parsed into a prompt.Template struct.
	dummyTemplate := `template: |
  Hello, {{ audience }}! You are learning about {{ topic }}.`
	writeFile(t, templatePath, dummyTemplate)

	// Create a dummy config file with valid fields.
	dummyConfig := `
audience: "Test Audience"
learning_stage: "beginner"
topic: "Test Topic"
context: "Test Context"
analogies: "Test Analogies"
concepts:
  - "Test Concept"
explanation_requirements:
  - "Test Requirement"
formatting:
  - "Test Formatting"
constraints:
  - "Test Constraint"
output_format:
  - "Test Output"
purpose: "Test Purpose"
tone: "Test Tone"
`
	writeFile(t, configPath, dummyConfig)

	// Capture output
	var out bytes.Buffer
	runner := &PromptRunner{Out: &out}

	// Act
	runner.Run()

	// Assert
	output := out.String()
	assert.Contains(t, output, "Prompt saved to: "+outputPath)

	// Validate that the generated prompt file contains the expected rendered output.
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Hello, Test Audience! You are learning about Test Topic.")
}

// Helper: writeFile creates a temporary file with the given content.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
}
