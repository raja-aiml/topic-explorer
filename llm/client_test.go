package llm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	llmConfig "raja.aiml/ai.explorer/config/llm"
	"raja.aiml/ai.explorer/llm/wrapper"
)

// --- MockModel implements wrapper.Model ---
type MockModel struct {
	mock.Mock
}

func (m *MockModel) Call(ctx context.Context, prompt string, opts ...wrapper.CallOption) (string, error) {
	args := m.Called(ctx, prompt, opts)
	return args.String(0), args.Error(1)
}

func (m *MockModel) GenerateContent(ctx context.Context, messages []wrapper.MessageContent, opts ...wrapper.CallOption) (*wrapper.ContentResponse, error) {
	args := m.Called(ctx, messages, opts)
	return args.Get(0).(*wrapper.ContentResponse), args.Error(1)
}

// --- MockProvider implements wrapper.Provider ---
type MockProvider struct {
	model wrapper.Model
	err   error
}

func (p *MockProvider) Init(providerName, modelName string) (wrapper.Model, error) {
	return p.model, p.err
}

func TestNewClient_Success(t *testing.T) {
	mockModel := new(MockModel)
	mockProvider := &MockProvider{model: mockModel}

	cfg := llmConfig.Config{
		Provider: "openai",
		Model: llmConfig.ModelConfig{
			Name:        "gpt-4",
			Temperature: 0.7,
		},
		Client: llmConfig.ClientConfig{
			Timeout:        time.Second,
			VerboseLogging: false,
		},
	}

	client, err := NewClient(cfg, mockProvider, func(ctx context.Context, model wrapper.Model, prompt string, opts ...wrapper.CallOption) (string, error) {
		return "mock-response", nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClient_Error(t *testing.T) {
	mockProvider := &MockProvider{err: errors.New("init error")}

	cfg := llmConfig.Config{
		Provider: "openai",
		Model:    llmConfig.ModelConfig{Name: "test"},
	}

	client, err := NewClient(cfg, mockProvider, nil)

	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to initialize LLM provider")
}

func TestClient_Chat_Success(t *testing.T) {
	mockModel := new(MockModel)
	mockModel.On("Call", mock.Anything, "test prompt", mock.Anything).Return("mock reply", nil)

	cfg := llmConfig.Config{
		Model: llmConfig.ModelConfig{
			Name:        "test",
			Temperature: 0.9,
		},
		Client: llmConfig.ClientConfig{
			Timeout:        2 * time.Second,
			VerboseLogging: false,
		},
	}

	client := &Client{
		model:  mockModel,
		config: cfg,
		callGen: func(ctx context.Context, model wrapper.Model, prompt string, opts ...wrapper.CallOption) (string, error) {
			return model.Call(ctx, prompt, opts...)
		},
	}

	resp, err := client.Chat(context.Background(), "test prompt")
	assert.NoError(t, err)
	assert.Equal(t, "mock reply", resp)
	mockModel.AssertExpectations(t)
}

func TestClient_Chat_Error(t *testing.T) {
	mockModel := new(MockModel)
	mockModel.On("Call", mock.Anything, "fail", mock.Anything).Return("", errors.New("chat failure"))

	cfg := llmConfig.Config{
		Model: llmConfig.ModelConfig{
			Name:        "test",
			Temperature: 0.5,
		},
		Client: llmConfig.ClientConfig{
			Timeout:        time.Second,
			VerboseLogging: true,
		},
	}

	client := &Client{
		model:  mockModel,
		config: cfg,
		callGen: func(ctx context.Context, model wrapper.Model, prompt string, opts ...wrapper.CallOption) (string, error) {
			return model.Call(ctx, prompt, opts...)
		},
	}

	resp, err := client.Chat(context.Background(), "fail")
	assert.Error(t, err)
	assert.Empty(t, resp)
	assert.Contains(t, err.Error(), "chat failed")
}
