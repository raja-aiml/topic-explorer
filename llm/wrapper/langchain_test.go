package wrapper_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
	"raja.aiml/ai.explorer/llm/wrapper"
)

// --- Mock Implementations ---

type mockLLM struct{ mock.Mock }

func (m *mockLLM) GenerateContent(ctx context.Context, msgs []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	args := m.Called(ctx, msgs, opts)
	return args.Get(0).(*llms.ContentResponse), args.Error(1)
}

func (m *mockLLM) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	args := m.Called(ctx, prompt, opts)
	return args.String(0), args.Error(1)
}

type mockEmbeddings struct{ mock.Mock }

func (m *mockEmbeddings) EmbedDocuments(ctx context.Context, inputs []string) ([][]float32, error) {
	args := m.Called(ctx, inputs)
	return args.Get(0).([][]float32), args.Error(1)
}

func (m *mockEmbeddings) EmbedQuery(ctx context.Context, input string) ([]float32, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]float32), args.Error(1)
}

// --- Tests ---

func TestProvider_Init(t *testing.T) {
	cases := []struct {
		name   string
		p, m   string
		hasErr bool
	}{
		{"openai", "openai", "gpt-3.5", false},
		{"ollama", "ollama", "phi4", false},
		{"invalid", "fake", "none", true},
	}

	prov := &wrapper.LangchaingoProvider{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.p == "openai" {
				_ = os.Setenv("OPENAI_API_KEY", "dummy-key")
				t.Cleanup(func() {
					_ = os.Unsetenv("OPENAI_API_KEY")
				})
			}

			model, err := prov.Init(c.p, c.m)
			if c.hasErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
			}
		})
	}
}

func TestGenerateFromSinglePrompt(t *testing.T) {
	ctx := context.Background()
	prompt := "test"

	mockModel := new(mockLLM)
	mockModel.On("GenerateContent", ctx, mock.Anything, mock.Anything).
		Return(newMockContentResponse("ok"), nil)

	resp, err := wrapper.GenerateFromSinglePrompt(ctx, mockModel, prompt)

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
	mockModel.AssertExpectations(t)
}

func TestEmbedderImpl_Embed(t *testing.T) {
	ctx := context.Background()
	inputs := []string{"foo", "bar"}
	expected := [][]float32{{1.0}, {2.0}}

	mockE := new(mockEmbeddings)
	mockE.On("EmbedDocuments", ctx, inputs).Return(expected, nil)

	e := wrapper.NewEmbedderFromBase(mockE)
	out, err := e.Embed(ctx, inputs)

	assert.NoError(t, err)
	assert.Equal(t, expected, out)
	mockE.AssertExpectations(t)
}

func TestWithTemperature(t *testing.T) {
	opt := wrapper.WithTemperature(0.5)
	assert.NotNil(t, opt)
}

func TestWithStreamingFunc(t *testing.T) {
	f := func(ctx context.Context, chunk []byte) error {
		if strings.Contains(string(chunk), "fail") {
			return errors.New("fail")
		}
		return nil
	}
	opt := wrapper.WithStreamingFunc(f)
	assert.NotNil(t, opt)
}

func newMockContentResponse(text string) *llms.ContentResponse {
	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{Content: text},
		},
	}
}
