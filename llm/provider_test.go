package llm

import (
	"fmt"
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
	model, err := initLLMProvider(cfg)
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
	model, err := initLLMProvider(cfg)
	if err != nil {
		t.Fatalf("Expected no error for provider 'openai', got: %v", err)
	}
	if model == nil {
		t.Fatal("Expected non-nil model for provider 'openai'")
	}
}

func TestInitLLMProvider_Unsupported(t *testing.T) {
	// Test the unsupported provider branch.
	cfg := llmConfig.Config{
		Provider: "unknown",
		Model: llmConfig.ModelConfig{
			Name: "dummy",
		},
	}
	model, err := initLLMProvider(cfg)
	if err == nil {
		t.Fatal("Expected error for unsupported provider, got nil")
	}
	if model != nil {
		t.Fatalf("Expected nil model for unsupported provider, got: %v", model)
	}
	expected := fmt.Sprintf("unsupported LLM provider: %s", cfg.Provider)
	if err.Error() != expected {
		t.Fatalf("Expected error message %q, got %q", expected, err.Error())
	}
}
