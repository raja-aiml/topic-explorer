package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/paths"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Generate a prompt and get an LLM response",
	Run: func(cmd *cobra.Command, args []string) {
		// Get paths from path manager
		topic = strings.ToLower(topic)

		// Get paths from path manager
		configPath = paths.GetConfigPath(topic, configPath)
		templatePath = paths.GetTemplatePath(templatePath)
		outputPath = paths.GetOutputPath(topic, outputPath)
		responseFilePath = paths.GetAnswerPath(topic, responseFilePath)
		runChat()
	},
}

// Init function to add chatCmd to root
func init() {
	// Define CLI flags
	chatCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic name for prompt generation (required)")
	chatCmd.Flags().StringVarP(&providerName, "provider", "p", DefaultProvider, "Provider name [ollama, openai] (optional)")
	chatCmd.Flags().StringVarP(&modelName, "model", "m", DefaultModel, "Model name [phi4, gpt-4o] (optional)")
	chatCmd.Flags().StringVarP(&outputPath, "promptOutput", "o", DefaultPromptPath, "Path to save generated prompt (optional)")

	chatCmd.MarkFlagRequired("topic") // Ensure topic is always required

	// Register command with root
	rootCmd.AddCommand(chatCmd)
}

// runChat first generates a prompt, then gets an LLM response
func runChat() {
	fmt.Println("Generating prompt...")
	promptFilePath := buildPrompt(templatePath, configPath, outputPath) // Calls buildPrompt from prompt.go

	fmt.Printf("Generated prompt saved to: %s\n", promptFilePath)

	fmt.Println("Reading generated prompt...")
	promptText, err := getPrompt(promptFilePath) // Calls getPrompt from llm.go
	if err != nil {
		log.Fatalf("Error retrieving prompt: %v", err)
	}

	fmt.Println("Generating LLM response...")

	response, err := runLLMInteraction(promptText) // Calls runLLMInteraction from llm.go
	if err != nil {
		log.Fatalf("Error generating LLM response: %v", err)
	}
	// response, err := client.Chat(context.Background(), promptText)

	fmt.Println("\nLLM Response:\n", response)
	// Check if topic is provided, if not, skip saving the response
	if topic == "" {
		return
	} else {
		fmt.Printf("Saving response to: %s\n", responseFilePath)
		saveResponse(response, responseFilePath) // Calls saveResponse from prompt.go
	}

}
