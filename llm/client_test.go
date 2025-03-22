package llm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

// MockModel mocks the llms.Model interface.
type MockModel struct {
	mock.Mock
}

func (m *MockModel) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	args := m.Called(ctx, prompt, opts)
	return args.String(0), args.Error(1)
}

func (m *MockModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	args := m.Called(ctx, messages, opts)
	return args.Get(0).(*llms.ContentResponse), args.Error(1)
}

func TestClient_Chat_Success(t *testing.T) {
	mockModel := new(MockModel)
	mockModel.On("Call", mock.Anything, "hello", mock.Anything).Return("mocked reply", nil)

	mockProvider := func(_ Config) (llms.Model, error) {
		return mockModel, nil
	}

	mockGenerator := func(ctx context.Context, model llms.Model, prompt string, opts ...llms.CallOption) (string, error) {
		return model.Call(ctx, prompt, opts...)
	}

	cfg := Config{
		Model: ModelConfig{
			Name:        "test",
			Temperature: 0.9,
		},
		Client: ClientConfig{
			Timeout:        time.Second,
			VerboseLogging: false,
		},
	}

	client, err := NewClient(cfg, mockProvider, mockGenerator)
	assert.NoError(t, err)

	resp, err := client.Chat(context.Background(), "hello")
	assert.NoError(t, err)
	assert.Equal(t, "mocked reply", resp)

	mockModel.AssertExpectations(t)
}
