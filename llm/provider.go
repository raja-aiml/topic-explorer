package llm

import (
	"fmt"

	"raja.aiml/ai.explorer/config/llm"
	langchainWrapper "raja.aiml/ai.explorer/llm/wrapper"
)

// initLLMProvider initializes the LLM model using the langchaingo wrapper.
func InitLLMProvider(cfg llm.Config) (langchainWrapper.Model, error) {
	provider := &langchainWrapper.LangchaingoProvider{}
	model, err := provider.Init(cfg.Provider, cfg.Model.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM provider: %w", err)
	}
	return model, nil
}
