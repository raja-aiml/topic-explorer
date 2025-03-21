package llm

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func TestDefaultStreamHandler(t *testing.T) {
	// Create a pipe to capture stdout.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Save original os.Stdout and ensure restoration after the test.
	origStdout := os.Stdout
	os.Stdout = w

	// Define test input.
	testChunk := []byte("Hello, test!")

	// Call the function under test.
	err = defaultStreamHandler(context.Background(), testChunk)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Close writer to flush.
	w.Close()

	// Capture the output.
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	// Restore original stdout.
	os.Stdout = origStdout

	// Compare output.
	output := buf.String()
	expected := "Hello, test!"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}
