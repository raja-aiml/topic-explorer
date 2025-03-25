package llm

import (
	"context"
	"errors"
	"math"

	"raja.aiml/ai.explorer/llm/wrapper"
)

// SimilarityService computes cosine similarity using an Embedder.
type SimilarityService struct {
	embedder wrapper.Embedder
}

// NewSimilarityService constructs a similarity service using the provided Embedder.
func NewSimilarityService(embedder wrapper.Embedder) *SimilarityService {
	return &SimilarityService{embedder: embedder}
}

// GetEmbeddings returns the embeddings for a slice of input strings.
// It is useful for retrieving raw vector representations for downstream tasks like storage or clustering.
func (s *SimilarityService) GetEmbeddings(ctx context.Context, inputs []string) ([][]float32, error) {
	if len(inputs) == 0 {
		return nil, errors.New("input list is empty")
	}
	return s.embedder.Embed(ctx, inputs)
}

// Compare computes cosine similarity between the embeddings of two input strings.
func (s *SimilarityService) Compare(ctx context.Context, a, b string) (float64, error) {
	vecs, err := s.embedder.Embed(ctx, []string{a, b})
	if err != nil {
		return 0, err
	}
	if len(vecs) < 2 {
		return 0, errors.New("not enough embeddings returned")
	}
	return cosine(vecs[0], vecs[1]), nil
}

// cosine calculates cosine similarity between two vectors.
func cosine(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
