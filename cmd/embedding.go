package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"raja.aiml/ai.explorer/llm"
	"raja.aiml/ai.explorer/llm/wrapper"
)

// EmbeddingRunner is responsible for executing embedding subcommands.
type EmbeddingRunner struct {
	Out io.Writer
}

func (r *EmbeddingRunner) RunGet(inputs []string) {
	if len(inputs) == 0 {
		log.Fatal("At least one input string is required")
	}

	ctx := context.Background()
	embedder, err := wrapper.NewOpenAIEmbedder()
	if err != nil {
		log.Fatalf("Failed to initialize embedder: %v", err)
	}
	service := llm.NewSimilarityService(embedder)

	fmt.Fprintln(r.Out, "Generating embeddings...")
	embeddings, err := service.GetEmbeddings(ctx, inputs)
	if err != nil {
		log.Fatalf("Error generating embeddings: %v", err)
	}

	for i, emb := range embeddings {
		fmt.Fprintf(r.Out, "Input: %q\nEmbedding: %v\n\n", inputs[i], emb)
	}
}

func (r *EmbeddingRunner) RunCompare(a, b string) {
	ctx := context.Background()
	embedder, err := wrapper.NewOpenAIEmbedder()
	if err != nil {
		log.Fatalf("Failed to initialize embedder: %v", err)
	}
	service := llm.NewSimilarityService(embedder)

	fmt.Fprintln(r.Out, "Computing cosine similarity...")
	score, err := service.Compare(ctx, a, b)
	if err != nil {
		log.Fatalf("Error comparing embeddings: %v", err)
	}

	fmt.Fprintf(r.Out, "Cosine similarity between:\n%q\nand\n%q\n→ %.4f\n", a, b, score)
}

var embeddingCmd = &cobra.Command{
	Use:   "embedding",
	Short: "Work with text embeddings and similarity",
}

var embeddingGetCmd = &cobra.Command{
	Use:   "get [text1 text2 ...]",
	Short: "Generate and print embeddings for input text",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runner := &EmbeddingRunner{Out: os.Stdout}
		runner.RunGet(args)
	},
}

var embeddingCompareCmd = &cobra.Command{
	Use:   "compare [text1] [text2]",
	Short: "Compute cosine similarity between two inputs",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runner := &EmbeddingRunner{Out: os.Stdout}
		runner.RunCompare(args[0], args[1])
	},
}

func init() {
	embeddingCmd.AddCommand(embeddingGetCmd)
	embeddingCmd.AddCommand(embeddingCompareCmd)
	rootCmd.AddCommand(embeddingCmd)
}
