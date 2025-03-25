package llm

import (
	"context"
	"errors"
	"math"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
)

// SimilarityService computes cosine similarity using an Embedder.
type SimilarityService struct {
	embedder embeddings.Embedder
}

// NewSimilarityService returns a service with a langchaingo-compatible Embedder.
func NewSimilarityService(embedder embeddings.Embedder) *SimilarityService {
	return &SimilarityService{embedder: embedder}
}

// Compare returns the cosine similarity between two input strings.
func (s *SimilarityService) Compare(ctx context.Context, a, b string) (float64, error) {
	vecs, err := s.embedder.EmbedDocuments(ctx, []string{a, b})
	if err != nil {
		return 0, err
	}
	if len(vecs) < 2 {
		return 0, errors.New("not enough embeddings returned")
	}
	return cosine(vecs[0], vecs[1]), nil
}

// cosine computes cosine similarity between two float32 vectors.
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

// NewOpenAIEmbedder uses langchaingo's default .env-based API key loading.
func NewOpenAIEmbedder() (embeddings.Embedder, error) {
	llm, err := openai.New() // auto-loads OPENAI_API_KEY from env
	if err != nil {
		return nil, err
	}
	return embeddings.NewEmbedder(llm)
}
