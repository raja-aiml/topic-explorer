package llm

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

// MockModel is a Testify mock that implements the llms.Model interface.
type MockModel struct {
	mock.Mock
}

// Call is a mock implementation of the Call method.
// It flattens the variadic opts so that expectations can be set easily.
func (m *MockModel) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	allArgs := []interface{}{ctx, prompt}
	for _, opt := range opts {
		allArgs = append(allArgs, opt)
	}
	// Dump the arguments for debugging.
	spew.Dump("MockModel.Call invoked with:", allArgs)
	args := m.Called(allArgs...)
	return args.String(0), args.Error(1)
}

// GenerateContent is a mock implementation of the GenerateContent method.
// Although it is not used in this test (because we override generateFromSinglePrompt),
// it must be implemented to satisfy the interface.
func (m *MockModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	allArgs := []interface{}{ctx, messages}
	for _, opt := range opts {
		allArgs = append(allArgs, opt)
	}
	spew.Dump("MockModel.GenerateContent invoked with:", allArgs)
	args := m.Called(allArgs...)
	return args.Get(0).(*llms.ContentResponse), args.Error(1)
}

func TestChat(t *testing.T) {
	mockModel := new(MockModel)
	// Expect exactly four arguments: context, "test prompt", and two call options.
	mockModel.
		On("Call", mock.Anything, "test prompt", mock.Anything, mock.Anything).
		Return("mocked response", nil)

	// Override generateFromSinglePrompt so that it calls model.Call.
	origGenerate := generateFromSinglePrompt
	generateFromSinglePrompt = func(ctx context.Context, model llms.Model, prompt string, opts ...llms.CallOption) (string, error) {
		return model.Call(ctx, prompt, opts...)
	}
	defer func() { generateFromSinglePrompt = origGenerate }()

	// Override initProvider to return our mock.
	origInitProvider := initProvider
	initProvider = func(cfg Config) (llms.Model, error) {
		return mockModel, nil
	}
	defer func() { initProvider = origInitProvider }()

	client, err := NewClient(Config{
		Provider: "dummy",
		Model: ModelConfig{
			Name:        "dummy-model",
			Temperature: 0.9,
		},
		Client: ClientConfig{
			Timeout:        1 * time.Second,
			VerboseLogging: true,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.Chat(context.Background(), "test prompt")
	if err != nil {
		t.Fatalf("Failed to chat: %v", err)
	}
	if response != "mocked response" {
		t.Fatalf("Expected response to be 'mocked response', but got '%s'", response)
	}
	mockModel.AssertExpectations(t)
}
