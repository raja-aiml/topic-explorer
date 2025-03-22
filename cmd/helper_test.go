package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getPrompt_success(t *testing.T) {
	content := "Hello Prompt"
	tmpFile := writeTempFile(t, content)

	result, err := getPrompt(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, content, result)
}

func Test_getPrompt_failure(t *testing.T) {
	_, err := getPrompt("nonexistent.txt")
	assert.Error(t, err)
}

func Test_saveResponse_success(t *testing.T) {
	tmpFile := t.TempDir() + "/response.txt"
	text := "LLM response"

	err := saveResponse(text, tmpFile)
	require.NoError(t, err)

	bytes, _ := os.ReadFile(tmpFile)
	assert.Equal(t, text, string(bytes))
}

func Test_saveResponse_failure(t *testing.T) {
	err := saveResponse("data", "/badpath/response.txt")
	assert.Error(t, err)
}

// Helper: write a temporary file
func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile := t.TempDir() + "/prompt.txt"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	return tmpFile
}
