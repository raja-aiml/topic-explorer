package llm

import (
	"context"
	"fmt"

	llmConfig "raja.aiml/ai.explorer/config/llm"
	"raja.aiml/ai.explorer/llm/wrapper"
)

// LLM defines the interface for any chat-capable client.
type LLM interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// Client wraps an LLM model and config.
type Client struct {
	model   wrapper.Model
	config  llmConfig.Config
	callGen func(ctx context.Context, model wrapper.Model, prompt string, opts ...wrapper.CallOption) (string, error)
}

// NewClient supports injecting dependencies for testability.
func NewClient(cfg llmConfig.Config, provider wrapper.Provider, generator func(context.Context, wrapper.Model, string, ...wrapper.CallOption) (string, error)) (*Client, error) {
	model, err := provider.Init(cfg.Provider, cfg.Model.Name)
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
	return NewClient(cfg, &wrapper.LangchaingoProvider{}, wrapper.GenerateFromSinglePrompt)
}

// Chat generates a response for the given prompt.
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Client.Timeout)
	defer cancel()

	opts := []wrapper.CallOption{
		wrapper.WithTemperature(c.config.Model.Temperature),
	}
	if c.config.Client.VerboseLogging {
		opts = append(opts, wrapper.WithStreamingFunc(defaultStreamHandler))
	}

	response, err := c.callGen(ctx, c.model, prompt, opts...)
	if err != nil {
		return "", fmt.Errorf("chat failed: %w", err)
	}
	return response, nil
}
