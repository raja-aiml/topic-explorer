package llm

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// mockEmbedder is a fake implementation of the Embedder interface for testing.
type mockEmbedder struct {
	output [][]float32
	err    error
}

func (m *mockEmbedder) Embed(_ context.Context, _ []string) ([][]float32, error) {
	return m.output, m.err
}

func TestGetEmbeddings(t *testing.T) {
	mock := &mockEmbedder{
		output: [][]float32{
			{0.1, 0.2, 0.3},
			{0.4, 0.5, 0.6},
		},
	}
	service := NewSimilarityService(mock)

	embeddings, err := service.GetEmbeddings(context.Background(), []string{"one", "two"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := mock.output
	if !reflect.DeepEqual(embeddings, expected) {
		t.Errorf("expected %v, got %v", expected, embeddings)
	}
}

func TestGetEmbeddings_EmptyInput(t *testing.T) {
	service := NewSimilarityService(&mockEmbedder{})

	_, err := service.GetEmbeddings(context.Background(), []string{})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestCompare(t *testing.T) {
	// cosine similarity between [1, 0] and [0, 1] is 0
	mock := &mockEmbedder{
		output: [][]float32{
			{1.0, 0.0},
			{0.0, 1.0},
		},
	}
	service := NewSimilarityService(mock)

	similarity, err := service.Compare(context.Background(), "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if similarity != 0 {
		t.Errorf("expected similarity 0, got %v", similarity)
	}
}

func TestCompare_ErrorFromEmbedder(t *testing.T) {
	mock := &mockEmbedder{
		err: errors.New("embedding failure"),
	}
	service := NewSimilarityService(mock)

	_, err := service.Compare(context.Background(), "a", "b")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCompare_NotEnoughEmbeddings(t *testing.T) {
	mock := &mockEmbedder{
		output: [][]float32{{1.0, 2.0}}, // only one embedding returned
	}
	service := NewSimilarityService(mock)

	_, err := service.Compare(context.Background(), "a", "b")
	if err == nil {
		t.Error("expected error for insufficient embeddings, got nil")
	}
}
