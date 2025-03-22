package llm

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

// ------------------------
// Mock Model Implementation
// ------------------------

type MockModel struct {
	mock.Mock
}

func (m *MockModel) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	allArgs := []interface{}{ctx, prompt}
	for _, opt := range opts {
		allArgs = append(allArgs, opt)
	}
	spew.Dump("MockModel.Call invoked with:", allArgs)
	args := m.Called(allArgs...)
	return args.String(0), args.Error(1)
}

func (m *MockModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	allArgs := []interface{}{ctx, messages}
	for _, opt := range opts {
		allArgs = append(allArgs, opt)
	}
	spew.Dump("MockModel.GenerateContent invoked with:", allArgs)
	args := m.Called(allArgs...)
	return args.Get(0).(*llms.ContentResponse), args.Error(1)
}

// ------------------------
// Helper Functions
// ------------------------

func setupMockModel() *MockModel {
	mockModel := new(MockModel)
	mockModel.On("Call", mock.Anything, "test prompt", mock.Anything, mock.Anything).
		Return("mocked response", nil)
	return mockModel
}

func overrideGenerateFromSinglePrompt() func() {
	original := generateFromSinglePrompt
	generateFromSinglePrompt = func(ctx context.Context, model llms.Model, prompt string, opts ...llms.CallOption) (string, error) {
		return model.Call(ctx, prompt, opts...)
	}
	return func() {
		generateFromSinglePrompt = original
	}
}

func overrideInitProvider(mockModel *MockModel) func() {
	original := initProvider
	initProvider = func(cfg Config) (llms.Model, error) {
		return mockModel, nil
	}
	return func() {
		initProvider = original
	}
}

func createTestClient(t *testing.T) *Client {
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
	return client
}

// ------------------------
// Main Test Function
// ------------------------

func TestChat(t *testing.T) {
	mockModel := setupMockModel()
	restoreGenerate := overrideGenerateFromSinglePrompt()
	defer restoreGenerate()
	restoreProvider := overrideInitProvider(mockModel)
	defer restoreProvider()

	client := createTestClient(t)

	response, err := client.Chat(context.Background(), "test prompt")
	if err != nil {
		t.Fatalf("Failed to chat: %v", err)
	}
	if response != "mocked response" {
		t.Fatalf("Expected response to be 'mocked response', but got '%s'", response)
	}

	mockModel.AssertExpectations(t)
}
