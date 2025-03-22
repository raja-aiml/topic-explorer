package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"raja.aiml/topic.explorer/llm"
	"raja.aiml/topic.explorer/paths"
)

// getPrompt reads the prompt from a file
func getPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading prompt file '%s': %w", path, err)
	}
	return string(content), nil
}

// saveResponse writes response to a file
func saveResponse(response, path string) error {
	return os.WriteFile(path, []byte(response), 0644)
}

// runLLMInteraction executes the LLM call
func runLLMInteraction(prompt string) (string, error) {
	cfg := llm.Config{
		Provider: providerName,
		Model:    llm.ModelConfig{Name: modelName, Temperature: temperature},
		Client:   llm.ClientConfig{Timeout: timeout, VerboseLogging: true},
	}
	client, err := llm.NewClient(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return client.Chat(ctx, prompt)
}

// resolvePaths sets common paths
func resolvePaths() {
	topic = strings.ToLower(topic)
	configPath = paths.GetConfigPath(topic, configPath)
	templatePath = paths.GetTemplatePath(templatePath)
	outputPath = paths.GetOutputPath(topic, outputPath)
	responseFilePath = paths.GetAnswerPath(topic, responseFilePath)
}
