package llm

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestConfigLoaderSuccess tests a successful YAML read and parse.
func TestConfigLoaderSuccess(t *testing.T) {
	// Save the original readFile function and restore it after the test.
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	// Create valid YAML content matching the Config structure.
	yamlContent := []byte(`
provider: "openai"
model:
  name: "gpt-4"
  temperature: 0.9
client:
  timeout: 120s
  verbose_logging: false
`)
	readFile = func(filename string) ([]byte, error) {
		return yamlContent, nil
	}

	cfg, err := ConfigLoader("dummy.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cfg.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got %q", cfg.Provider)
	}
	if cfg.Model.Name != "gpt-4" {
		t.Errorf("Expected model name 'gpt-4', got %q", cfg.Model.Name)
	}
	if cfg.Model.Temperature != 0.9 {
		t.Errorf("Expected temperature 0.9, got %f", cfg.Model.Temperature)
	}
	// "120s" should be parsed as 120 seconds.
	expectedTimeout := 120 * time.Second
	if cfg.Client.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, cfg.Client.Timeout)
	}
	if cfg.Client.VerboseLogging != false {
		t.Errorf("Expected verbose_logging false, got %v", cfg.Client.VerboseLogging)
	}
}

// TestConfigLoaderFileReadError simulates a file read error.
func TestConfigLoaderFileReadError(t *testing.T) {
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	expectedErr := errors.New("read error")
	readFile = func(filename string) ([]byte, error) {
		return nil, expectedErr
	}

	_, err := ConfigLoader("dummy.yaml")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected error to contain 'failed to read config file', got %v", err)
	}
	if !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected error to wrap underlying error %v, got %v", expectedErr, err)
	}
}

// TestConfigLoaderParseError simulates an error during YAML parsing.
func TestConfigLoaderParseError(t *testing.T) {
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	// Return invalid YAML content.
	readFile = func(filename string) ([]byte, error) {
		return []byte("invalid: [unclosed"), nil
	}

	_, err := ConfigLoader("dummy.yaml")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse config YAML") {
		t.Errorf("Expected error to contain 'failed to parse config YAML', got %v", err)
	}
}
