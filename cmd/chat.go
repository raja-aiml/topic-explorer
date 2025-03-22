package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/paths"
)

type ChatRunner struct {
	Out               io.Writer
	BuildPrompt       func(templatePath, configPath, outputPath string) string
	GetPrompt         func(promptFile string) (string, error)
	RunLLMInteraction func(prompt string) (string, error)
	SaveResponse      func(response, filePath string) error // ✅ corrected signature
}

func (cr *ChatRunner) Run() {
	fmt.Fprintln(cr.Out, "Generating prompt...")
	promptFilePath := cr.BuildPrompt(templatePath, configPath, outputPath)

	fmt.Fprintf(cr.Out, "Generated prompt saved to: %s\n", promptFilePath)

	fmt.Fprintln(cr.Out, "Reading generated prompt...")
	promptText, err := cr.GetPrompt(promptFilePath)
	if err != nil {
		log.Fatalf("Error retrieving prompt: %v", err)
	}

	fmt.Fprintln(cr.Out, "Generating LLM response...")
	response, err := cr.RunLLMInteraction(promptText)
	if err != nil {
		log.Fatalf("Error generating LLM response: %v", err)
	}

	fmt.Fprintf(cr.Out, "\nLLM Response:\n%s\n", response)

	if topic == "" {
		return
	}

	fmt.Fprintf(cr.Out, "Saving response to: %s\n", responseFilePath)
	if err := cr.SaveResponse(response, responseFilePath); err != nil {
		log.Fatalf("Error saving response: %v", err)
	}
}

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
