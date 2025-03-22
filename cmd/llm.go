package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/llm"
)

type LLMRunner struct {
	Out          io.Writer
	GetPrompt    func(promptPath string) (string, error)
	RunLLM       func(prompt string) (string, error)
	SaveResponse func(response, responseFilePath string) error
}

func (r *LLMRunner) Run() {
	fmt.Fprintln(r.Out, "Reading prompt...")
	promptText, err := r.GetPrompt(promptPath)
	if err != nil {
		log.Fatalf("Failed to retrieve prompt: %v", err)
	}

	fmt.Fprintln(r.Out, "Generating response...")
	response, err := r.RunLLM(promptText)
	if err != nil {
		log.Fatalf("Failed to generate LLM response: %v", err)
	}

	fmt.Fprintf(r.Out, "\nLLM Response:\n%s\n", response)

	if responseFilePath != "" {
		fmt.Fprintf(r.Out, "Saving response to: %s\n", responseFilePath)
		if err := r.SaveResponse(response, responseFilePath); err != nil {
			log.Fatalf("Failed to save response: %v", err)
		}
	}
}

// llmCmd represents the LLM command
var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Interact with the LLM model",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup runner
		runner := &LLMRunner{
			Out:          os.Stdout,
			GetPrompt:    getPrompt,
			RunLLM:       runLLMInteraction,
			SaveResponse: saveResponse,
		}
		runner.Run()
	},
}

func init() {
	llmCmd.Flags().StringVarP(&providerName, "provider", "l", DefaultProvider, fmt.Sprintf("LLM provider to use, default is %s", DefaultProvider))
	llmCmd.Flags().StringVarP(&modelName, "model", "m", DefaultModel, fmt.Sprintf("LLM model to use, default is %s", DefaultModel))
	llmCmd.Flags().Float64VarP(&temperature, "temperature", "t", DefaultTemperature, fmt.Sprintf("Sampling temperature, default is %.1f", DefaultTemperature))
	llmCmd.Flags().StringVarP(&promptPath, "prompt", "p", DefaultPromptPath, fmt.Sprintf("Path to prompt file, default: %s", DefaultPromptPath))
	llmCmd.Flags().DurationVarP(&timeout, "timeout", "d", DefaultTimeout, fmt.Sprintf("Timeout for LLM request, default: %s", DefaultTimeout))
	llmCmd.Flags().StringVarP(&responseFilePath, "save", "s", "", "Optional path to save LLM response")

	rootCmd.AddCommand(llmCmd)
}

// Helpers

func getPrompt(promptPath string) (string, error) {
	if promptPath == "" {
		promptPath = DefaultPromptPath
	}
	content, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("error reading prompt file '%s': %w", promptPath, err)
	}
	return string(content), nil
}

func runLLMInteraction(prompt string) (string, error) {
	fmt.Printf("Using Provider: %s | Model: %s\n", providerName, modelName)

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

	client, err := llm.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := client.Chat(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("chat failed: %w", err)
	}

	return response, nil
}

func saveResponse(response, filePath string) error {
	if err := os.WriteFile(filePath, []byte(response), 0644); err != nil {
		return fmt.Errorf("error saving response to '%s': %w", filePath, err)
	}
	return nil
}
