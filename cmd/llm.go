package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/llm"
)

// llmCmd represents the LLM command
var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Interact with the LLM model",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the prompt text (from file if provided, otherwise default file)
		finalPrompt, err := getPrompt(promptPath)
		if err != nil {
			log.Fatalf("Failed to retrieve prompt: %v", err)
		}

		// Run the LLM interaction with the resolved prompt
		runLLMInteraction(finalPrompt)
	},
}

// Init function to add llmCmd to root
func init() {
	// Define flags for customization with proper defaults
	llmCmd.Flags().StringVarP(&providerName, "provider", "l", DefaultProvider, fmt.Sprintf("LLM provider to use, default is %s", DefaultProvider))
	llmCmd.Flags().StringVarP(&modelName, "model", "m", DefaultModel, fmt.Sprintf("LLM model to use, default is %s", DefaultModel))
	llmCmd.Flags().Float64VarP(&temperature, "temperature", "t", DefaultTemperature, fmt.Sprintf("Sampling temperature, default is %.1f", DefaultTemperature))
	llmCmd.Flags().StringVarP(&promptPath, "prompt", "p", DefaultPromptPath, fmt.Sprintf("Path to a prompt file. If not provided, the default prompt file (%s) is used.", DefaultPromptPath))
	llmCmd.Flags().DurationVarP(&timeout, "timeout", "d", DefaultTimeout, fmt.Sprintf("Request timeout duration, default is %s", DefaultTimeout))

	// Register command with root
	rootCmd.AddCommand(llmCmd)
}

// getPrompt reads the prompt from a file (either user-specified or the default file).
func getPrompt(promptPath string) (string, error) {
	if promptPath == "" {
		promptPath = DefaultPromptPath // Use default prompt file if no file is provided
	}

	// Read from file using os.ReadFile (Go 1.16+ recommended method)
	content, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("error reading prompt file '%s': %w", promptPath, err)
	}

	return string(content), nil
}

// runLLMInteraction initializes the LLM client, generates a response, and returns it as a string
func runLLMInteraction(prompt string) (string, error) {
	fmt.Printf(" Provider: %s\n Model: %s\n: ", providerName, modelName)

	// Create LLM configuration
	config := llm.Config{
		Provider: providerName,
		Model: llm.ModelConfig{
			Name:        modelName,
			Temperature: temperature,
		},
		Client: llm.ClientConfig{
			Timeout:        timeout,
			VerboseLogging: true,
		},
	}

	// Initialize the LLM client
	client, err := llm.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Generate response
	response, err := client.Chat(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("chat failed: %w", err)
	}

	// Return the response as a string
	return response, nil
}

// saveResponse saves the LLM response to a file
func saveResponse(response, responseFilePath string) error {
	// Write to file using os.WriteFile (Go 1.16+ recommended method)
	err := os.WriteFile(responseFilePath, []byte(response), 0644)
	if err != nil {
		return fmt.Errorf("error writing response to file: %w", err)
	}

	return nil
}
