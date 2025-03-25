package wrapper

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

type (
	Model      = llms.Model
	CallOption = llms.CallOption
)

// ---------- LLM Provider Abstraction ----------

// Provider defines an abstraction to initialize LLM models.
type Provider interface {
	Init(providerName, modelName string) (Model, error)
}

// LangchaingoProvider is a concrete LLM provider using langchaingo.
type LangchaingoProvider struct{}

// Init returns a new Model for the given provider and model name.
func (p *LangchaingoProvider) Init(providerName, modelName string) (Model, error) {
	switch providerName {
	case "ollama":
		return ollama.New(ollama.WithModel(modelName))
	case "openai":
		return openai.New(openai.WithModel(modelName))
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", providerName)
	}
}

// ---------- Embedding Abstraction ----------

// Embedder defines a minimal interface for generating vector embeddings.
type Embedder interface {
	Embed(ctx context.Context, inputs []string) ([][]float32, error)
}

// embedderImpl wraps langchaingo's embeddings.Embedder.
type embedderImpl struct {
	base embeddings.Embedder
}

// NewOpenAIEmbedder creates an OpenAI-based Embedder (API key from env).
func NewOpenAIEmbedder() (Embedder, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, err
	}
	base, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	return &embedderImpl{base: base}, nil
}

// Embed returns the embeddings for the provided inputs.
func (e *embedderImpl) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	return e.base.EmbedDocuments(ctx, inputs)
}

// ---------- LLM Generation Helpers ----------

// GenerateFromSinglePrompt is an alias to langchaingo's llms.GenerateFromSinglePrompt
func GenerateFromSinglePrompt(ctx context.Context, model Model, prompt string, opts ...CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, model, prompt, opts...)
}

// WithTemperature wraps llms.WithTemperature
func WithTemperature(temp float64) CallOption {
	return llms.WithTemperature(temp)
}

// WithStreamingFunc wraps llms.WithStreamingFunc
func WithStreamingFunc(f func(ctx context.Context, chunk []byte) error) CallOption {
	return llms.WithStreamingFunc(f)
}
