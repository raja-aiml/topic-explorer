package cmd

import "time"

// Default file paths
const (
	DefaultConfigPath  = "resources/default/config.yaml" // Default config file
	DefaultProvider    = "ollama"
	DefaultModel       = "phi4"
	DefaultTemperature = 0.8
	DefaultTimeout     = 2 * time.Minute
	DefaultPromptPath  = "resources/default/prompt.txt" // Ensure a valid default prompt path
)

// Define flags as global variables to avoid scope issues
var (
	topic            string
	templatePath     string
	configPath       string
	outputPath       string
	responseFilePath string
)

// CLI flags
var (
	providerName string
	modelName    string
	temperature  float64
	promptPath   string
	timeout      time.Duration
)
