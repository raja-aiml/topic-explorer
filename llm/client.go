package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	llmConfig "raja.aiml/ai.explorer/config/llm"
)

// LLM defines the interface for any chat-capable client.
type LLM interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// Client wraps an LLM model and config.
type Client struct {
	model   llms.Model
	config  llmConfig.Config
	callGen func(ctx context.Context, model llms.Model, prompt string, opts ...llms.CallOption) (string, error)
}

// NewClient supports injecting dependencies for testability.
func NewClient(cfg llmConfig.Config, provider func(llmConfig.Config) (llms.Model, error), generator func(context.Context, llms.Model, string, ...llms.CallOption) (string, error)) (*Client, error) {
	model, err := provider(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM provider: %w", err)
	}
	return &Client{
		model:   model,
		config:  cfg,
		callGen: generator,
	}, nil
}

// NewDefaultClient returns a client with default dependencies.
func NewDefaultClient(cfg llmConfig.Config) (*Client, error) {
	return NewClient(cfg, initLLMProvider, llms.GenerateFromSinglePrompt)
}

// Chat generates a response for the given prompt.
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Client.Timeout)
	defer cancel()

	opts := []llms.CallOption{
		llms.WithTemperature(c.config.Model.Temperature),
	}
	if c.config.Client.VerboseLogging {
		opts = append(opts, llms.WithStreamingFunc(defaultStreamHandler))
	}

	response, err := c.callGen(ctx, c.model, prompt, opts...)
	if err != nil {
		return "", fmt.Errorf("chat failed: %w", err)
	}
	return response, nil
}
