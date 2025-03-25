package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"raja.aiml/ai.explorer/prompt"
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

	_ = promptCmd.MarkFlagRequired("topic")

	rootCmd.AddCommand(promptCmd)
}

// buildPrompt determines type from config and dispatches the right generator
func buildPrompt(tmpl, cfg, out string) string {
	switch detectPromptType(cfg) {
	case "topic":
		prompt.BuildTopicPrompt(tmpl, cfg, out)
	default:
		exitWithError(fmt.Errorf("unsupported prompt type detected from config"))
	}
	return out
}

// detectPromptType peeks into the YAML keys to infer prompt style
func detectPromptType(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		exitWithError(fmt.Errorf("failed to read config: %w", err))
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		exitWithError(fmt.Errorf("failed to parse config YAML: %w", err))
	}

	switch {
	case hasKey(raw, "audience"):
		return "topic"
	default:
		return "unknown"
	}
}

func hasKey(m map[string]any, key string) bool {
	_, ok := m[key]
	return ok
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
