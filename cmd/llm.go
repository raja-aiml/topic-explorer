package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type LLMRunner struct {
	Out          io.Writer
	GetPrompt    func(promptPath string) (string, error)
	RunLLM       func(prompt string) (string, error)
	SaveResponse func(response, path string) error
}

func (r *LLMRunner) Run() {
	fmt.Fprintln(r.Out, "Reading prompt...")
	text, err := r.GetPrompt(promptPath)
	if err != nil {
		log.Fatalf("Prompt error: %v", err)
	}

	fmt.Fprintln(r.Out, "Calling LLM...")
	resp, err := r.RunLLM(text)
	if err != nil {
		log.Fatalf("LLM error: %v", err)
	}

	fmt.Fprintf(r.Out, "\nLLM Response:\n%s\n", resp)

	if responseFilePath != "" {
		fmt.Fprintf(r.Out, "Saving to: %s\n", responseFilePath)
		if err := r.SaveResponse(resp, responseFilePath); err != nil {
			log.Fatalf("Save error: %v", err)
		}
	}
}

var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Send a raw prompt to LLM",
	Run: func(cmd *cobra.Command, args []string) {
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
	llmCmd.Flags().StringVarP(&providerName, "provider", "l", DefaultProvider, "LLM provider")
	llmCmd.Flags().StringVarP(&modelName, "model", "m", DefaultModel, "LLM model")
	llmCmd.Flags().Float64VarP(&temperature, "temperature", "t", DefaultTemperature, "Temperature")
	llmCmd.Flags().StringVarP(&promptPath, "prompt", "p", DefaultPromptPath, "Prompt file")
	llmCmd.Flags().DurationVarP(&timeout, "timeout", "d", DefaultTimeout, "Timeout duration")
	llmCmd.Flags().StringVarP(&responseFilePath, "save", "s", "", "Save response to file")
	rootCmd.AddCommand(llmCmd)
}
