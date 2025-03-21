package llm

import (
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// initLLMProvider initializes the LLM provider based on configuration.
func initLLMProvider(config Config) (llms.Model, error) {
	switch config.Provider {
	case "ollama":
		return ollama.New(ollama.WithModel(config.Model.Name))
	case "openai":
		return openai.New(openai.WithModel(config.Model.Name))
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", config.Provider)
	}
}
