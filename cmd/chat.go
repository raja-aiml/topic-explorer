package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type ChatRunner struct {
	Out io.Writer
}

func (r *ChatRunner) Run() {
	resolvePaths()

	fmt.Fprintln(r.Out, "Generating prompt...")
	promptFile := buildPrompt(templatePath, configPath, outputPath)

	fmt.Fprintln(r.Out, "Reading prompt...")
	text, err := getPrompt(promptFile)
	if err != nil {
		log.Fatalf("Prompt read error: %v", err)
	}

	fmt.Fprintln(r.Out, "Calling LLM...")
	resp, err := runLLMInteraction(text)
	if err != nil {
		log.Fatalf("LLM error: %v", err)
	}

	fmt.Fprintf(r.Out, "\nLLM Response:\n%s\n", resp)

	if topic != "" {
		fmt.Fprintf(r.Out, "Saving response to: %s\n", responseFilePath)
		if err := saveResponse(resp, responseFilePath); err != nil {
			log.Fatalf("Save error: %v", err)
		}
	}
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Generate a prompt and call LLM",
	Run: func(cmd *cobra.Command, args []string) {
		(&ChatRunner{Out: os.Stdout}).Run()
	},
}

func init() {
	chatCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic name (required)")
	chatCmd.Flags().StringVarP(&providerName, "provider", "p", DefaultProvider, "LLM provider")
	chatCmd.Flags().StringVarP(&modelName, "model", "m", DefaultModel, "LLM model name")
	chatCmd.Flags().StringVarP(&outputPath, "promptOutput", "o", DefaultPromptPath, "Prompt output path")
	chatCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(chatCmd)
}
