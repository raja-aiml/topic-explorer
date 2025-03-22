package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLLMRunnerRun(t *testing.T) {
	// Setup temporary prompt and output paths
	tmpDir := t.TempDir()
	promptPath = tmpDir + "/prompt.txt"
	responseFilePath = tmpDir + "/response.txt"

	expectedPrompt := "What is Go?"
	expectedResponse := "Go is an open-source programming language."

	// Create prompt file
	err := os.WriteFile(promptPath, []byte(expectedPrompt), 0644)
	require.NoError(t, err)

	// Capture CLI output
	var out bytes.Buffer

	// Create mock runner
	runner := &LLMRunner{
		Out: &out,
		GetPrompt: func(path string) (string, error) {
			assert.Equal(t, promptPath, path)
			return expectedPrompt, nil
		},
		RunLLM: func(prompt string) (string, error) {
			assert.Equal(t, expectedPrompt, prompt)
			return expectedResponse, nil
		},
		SaveResponse: func(response, path string) error {
			assert.Equal(t, expectedResponse, response)
			assert.Equal(t, responseFilePath, path)
			return os.WriteFile(path, []byte(response), 0644)
		},
	}

	// Act
	runner.Run()

	// Assert output
	output := out.String()
	assert.Contains(t, output, "Reading prompt...")
	assert.Contains(t, output, "Calling LLM...")
	assert.Contains(t, output, "LLM Response:")
	assert.Contains(t, output, "Saving to:")

	// Assert file content
	data, err := os.ReadFile(responseFilePath)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, string(data))
}
