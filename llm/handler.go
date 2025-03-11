package llm

import (
	"context"
	"fmt"
)

// StreamHandler defines a function signature for streaming chunk processing
type StreamHandler func(ctx context.Context, chunk []byte) error

// defaultStreamHandler prints output chunks to stdout
func defaultStreamHandler(ctx context.Context, chunk []byte) error {
	fmt.Print(string(chunk))
	return nil
}
