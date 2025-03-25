package llm

import (
	"errors"
	"os"
	"testing"

	llmConfig "raja.aiml/ai.explorer/config/llm"
)

func TestInitLLMProvider_Ollama(t *testing.T) {
	// Test the "ollama" branch.
	cfg := llmConfig.Config{
		Provider: "ollama",
		Model: llmConfig.ModelConfig{
			Name: "phi4",
		},
	}
	model, err := InitLLMProvider(cfg)
	if err != nil {
		t.Fatalf("Expected no error for provider 'ollama', got: %v", err)
	}
	if model == nil {
		t.Fatal("Expected non-nil model for provider 'ollama'")
	}
}

func TestInitLLMProvider_OpenAI(t *testing.T) {
	// Test the "openai" branch.
	os.Setenv("OPENAI_API_KEY", "dummy-key")
	defer os.Unsetenv("OPENAI_API_KEY")
	cfg := llmConfig.Config{
		Provider: "openai",
		Model: llmConfig.ModelConfig{
			Name: "gpt-4",
		},
	}
	model, err := InitLLMProvider(cfg)
	if err != nil {
		t.Fatalf("Expected no error for provider 'openai', got: %v", err)
	}
	if model == nil {
		t.Fatal("Expected non-nil model for provider 'openai'")
	}
}

func TestInitLLMProvider_Unsupported(t *testing.T) {
	cfg := llmConfig.Config{
		Provider: "unknown",
		Model: llmConfig.ModelConfig{
			Name: "dummy",
		},
	}
	model, err := InitLLMProvider(cfg)

	if err == nil {
		t.Fatal("Expected error for unsupported provider, got nil")
	}
	if model != nil {
		t.Fatalf("Expected nil model for unsupported provider, got: %v", model)
	}

	// Use errors.As to extract the custom error type if needed,
	// or just check for substring
	expected := "unsupported LLM provider: unknown"
	if err.Error() != expected && !containsError(err, expected) {
		t.Fatalf("Expected error to contain %q, got %q", expected, err.Error())
	}
}

// containsError is a helper that recursively unwraps errors
// and checks if the message contains the expected string.
func containsError(err error, expected string) bool {
	for err != nil {
		if err.Error() == expected {
			return true
		}
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			break
		}
		err = unwrapped
	}
	return false
}
