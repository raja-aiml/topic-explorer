package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"raja.aiml/ai.explorer/llm"
	"raja.aiml/ai.explorer/paths"

	llmConfig "raja.aiml/ai.explorer/config/llm"
)

// getPrompt reads the prompt from a file and returns its content as a string.
func getPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading prompt file '%s': %w", path, err)
	}
	return string(content), nil
}

// saveResponse writes the LLM response to the specified file.
func saveResponse(response, path string) error {
	return os.WriteFile(path, []byte(response), 0644)
}

// runLLMInteraction initializes the LLM client and returns the response for the given prompt.
func runLLMInteraction(prompt string) (string, error) {
	cfg := llmConfig.Config{
		Provider: providerName,
		Model: llmConfig.ModelConfig{
			Name:        modelName,
			Temperature: temperature,
		},
		Client: llmConfig.ClientConfig{
			Timeout:        timeout,
			VerboseLogging: true,
		},
	}

	client, err := llm.NewDefaultClient(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return client.Chat(ctx, prompt)
}

// resolvePaths fills in default or derived paths based on the topic and other flags.
func resolvePaths() {
	topic = strings.ToLower(topic)
	configPath = paths.GetConfigPath(topic, configPath)
	templatePath = paths.GetTemplatePath(templatePath)
	outputPath = paths.GetOutputPath(topic, outputPath)
	responseFilePath = paths.GetAnswerPath(topic, responseFilePath)
}
