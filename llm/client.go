package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

// Package-level variables for dependency injection.
// These can be overridden in tests.
var initProvider = initLLMProvider
var generateFromSinglePrompt = llms.GenerateFromSinglePrompt

// Client wraps LLM functionality providing a unified interface.
type Client struct {
	llm    llms.Model
	config Config
}

// NewClient initializes a new LLM client.
func NewClient(config Config) (*Client, error) {
	llm, err := initProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM provider: %w", err)
	}
	return &Client{
		llm:    llm,
		config: config,
	}, nil
}

// Chat generates a response for the given prompt.
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Client.Timeout)
	defer cancel()

	fmt.Printf("\n")
	fmt.Printf("Generating response for prompt (length: %d chars)", len(prompt))
	fmt.Printf("\n")
	opts := []llms.CallOption{
		llms.WithTemperature(c.config.Model.Temperature),
	}
	if c.config.Client.VerboseLogging {
		opts = append(opts, llms.WithStreamingFunc(defaultStreamHandler))
	}

	response, err := generateFromSinglePrompt(ctx, c.llm, prompt, opts...)
	if err != nil {
		fmt.Printf("Chat generation failed: %v", err)
		return "", fmt.Errorf("generation failed: %w", err)
	}
	fmt.Printf("\n")
	fmt.Printf("Chat generation completed successfully (length: %d chars)", len(response))
	return response, nil
}
