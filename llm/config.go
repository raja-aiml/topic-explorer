package llm

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Default configuration values
const (
	DefaultProvider       = "ollama"
	DefaultModelName      = "phi4"
	DefaultTemperature    = 0.8
	DefaultTimeout        = 2 * time.Minute
	DefaultVerboseLogging = true
)

// ModelConfig holds configuration specific to the language model
type ModelConfig struct {
	Name        string  // Name of the model
	Temperature float64 // Temperature setting
}

// ClientConfig holds runtime behavior configuration
type ClientConfig struct {
	Timeout        time.Duration // Maximum request time
	VerboseLogging bool          // Enable verbose logs
}

// Config aggregates model and client configurations
type Config struct {
	Provider string
	Model    ModelConfig
	Client   ClientConfig
}

// ConfigLoader loads LLM configuration from a YAML file
func ConfigLoader(filePath string) (Config, error) {
	config := Config{}

	// Read YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return config, nil
}
