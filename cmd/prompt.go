package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/prompt"
)

type PromptRunner struct {
	Out io.Writer
}

func (r *PromptRunner) Run() {
	resolvePaths()
	outPath := buildPrompt(templatePath, configPath, outputPath)
	fmt.Fprintf(r.Out, "Prompt saved to: %s\n", outPath)
}

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate prompt from YAML + config",
	Run: func(cmd *cobra.Command, args []string) {
		(&PromptRunner{Out: os.Stdout}).Run()
	},
}

func init() {
	promptCmd.Flags().StringVarP(&topic, "topic", "", "", "Topic name (required)")
	promptCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to template YAML")
	promptCmd.Flags().StringVarP(&configPath, "config", "c", "", "Config YAML path")
	promptCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Generated prompt output path")
	promptCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(promptCmd)
}

func buildPrompt(tmpl, cfg, out string) string {
	prompt.Build(tmpl, cfg, out)
	return out
}
