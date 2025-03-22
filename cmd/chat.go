package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/paths"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Generate a prompt and get an LLM response",
	Run: func(cmd *cobra.Command, args []string) {
		topic = strings.ToLower(topic)

		// Resolve paths
		configPath = paths.GetConfigPath(topic, configPath)
		templatePath = paths.GetTemplatePath(templatePath)
		outputPath = paths.GetOutputPath(topic, outputPath)
		responseFilePath = paths.GetAnswerPath(topic, responseFilePath)

		// Create a runner with real implementations
		runner := &ChatRunner{
			Out:               os.Stdout,
			BuildPrompt:       buildPrompt,
			GetPrompt:         getPrompt,
			RunLLMInteraction: runLLMInteraction,
			SaveResponse:      saveResponse,
		}
		runner.Run()
	},
}

// Init adds chat command to the root
func init() {
	chatCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic name for prompt generation (required)")
	chatCmd.Flags().StringVarP(&providerName, "provider", "p", DefaultProvider, "Provider name [ollama, openai] (optional)")
	chatCmd.Flags().StringVarP(&modelName, "model", "m", DefaultModel, "Model name [phi4, gpt-4o] (optional)")
	chatCmd.Flags().StringVarP(&outputPath, "promptOutput", "o", DefaultPromptPath, "Path to save generated prompt (optional)")

	chatCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(chatCmd)
}
