// package llm

// import (
// 	"context"
// 	"fmt"

// 	"github.com/tmc/langchaingo/llms"
// )

// type LLM interface {
// 	Chat(ctx context.Context, prompt string) (string, error)
// }

// // Package-level variables for dependency injection.
// // These can be overridden in tests.
// var initProvider = initLLMProvider
// var generateFromSinglePrompt = llms.GenerateFromSinglePrompt
// var NewClient = newClient // make NewClient assignable

// // Client wraps LLM functionality providing a unified interface.
// type Client struct {
// 	llm    llms.Model
// 	config Config
// }

// // NewClient initializes a new LLM client.
// func newClient(config Config) (*Client, error) {
// 	llm, err := initProvider(config)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize LLM provider: %w", err)
// 	}
// 	return &Client{
// 		llm:    llm,
// 		config: config,
// 	}, nil
// }

// // Chat generates a response for the given prompt.
// func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
// 	ctx, cancel := context.WithTimeout(ctx, c.config.Client.Timeout)
// 	defer cancel()

// 	fmt.Printf("\n")
// 	fmt.Printf("Generating response for prompt (length: %d chars)", len(prompt))
// 	fmt.Printf("\n")
// 	opts := []llms.CallOption{
// 		llms.WithTemperature(c.config.Model.Temperature),
// 	}
// 	if c.config.Client.VerboseLogging {
// 		opts = append(opts, llms.WithStreamingFunc(defaultStreamHandler))
// 	}

// 	response, err := generateFromSinglePrompt(ctx, c.llm, prompt, opts...)
// 	if err != nil {
// 		fmt.Printf("Chat generation failed: %v", err)
// 		return "", fmt.Errorf("generation failed: %w", err)
// 	}
// 	fmt.Printf("\n")
// 	fmt.Printf("Chat generation completed successfully (length: %d chars)", len(response))
// 	return response, nil
// }

package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

// LLM defines the interface for any chat-capable client.
type LLM interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// Client wraps an LLM model and config.
type Client struct {
	model   llms.Model
	config  Config
	callGen func(ctx context.Context, model llms.Model, prompt string, opts ...llms.CallOption) (string, error)
}

// NewClient supports injecting dependencies for testability.
func NewClient(cfg Config, provider func(Config) (llms.Model, error), generator func(context.Context, llms.Model, string, ...llms.CallOption) (string, error)) (*Client, error) {
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
func NewDefaultClient(cfg Config) (*Client, error) {
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
