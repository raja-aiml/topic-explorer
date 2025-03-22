package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/paths"
	"raja.aiml/topic.explorer/prompt"
)

type PromptRunner struct {
	Out         io.Writer
	BuildPrompt func(templatePath, configPath, outputPath string) string
}

func (r *PromptRunner) Run() {
	topic = strings.ToLower(topic)

	configPath = paths.GetConfigPath(topic, configPath)
	templatePath = paths.GetTemplatePath(templatePath)
	outputPath = paths.GetOutputPath(topic, outputPath)

	generatedPath := r.BuildPrompt(templatePath, configPath, outputPath)

	fmt.Fprintf(r.Out, "Prompt generated and saved to: %s\n", generatedPath)
}

// promptCmd represents the prompt generation command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate a structured prompt from YAML templates",
	Run: func(cmd *cobra.Command, args []string) {
		runner := &PromptRunner{
			Out:         os.Stdout,
			BuildPrompt: buildPrompt,
		}
		runner.Run()
	},
}

func init() {
	promptCmd.Flags().StringVarP(&topic, "topic", "", "", "Topic name (required)")
	promptCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to template file (default: resources/template.yaml)")
	promptCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file (default: based on topic)")
	promptCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to save generated prompt (default: based on topic)")

	promptCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(promptCmd)
}

// buildPrompt generates the prompt and returns the output file path
func buildPrompt(templatePath, configPath, outputPath string) string {
	prompt.Build(templatePath, configPath, outputPath)
	return outputPath
}
