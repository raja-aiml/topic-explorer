package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"raja.aiml/topic.explorer/paths"
	"raja.aiml/topic.explorer/prompt"
)

// promptCmd represents the prompt generation command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate a structured prompt from YAML templates",
	Run: func(cmd *cobra.Command, args []string) {
		topic = strings.ToLower(topic)

		// Get paths from path manager
		configPath = paths.GetConfigPath(topic, configPath)
		outputPath = paths.GetOutputPath(topic, outputPath)
		templatePath = paths.GetTemplatePath(templatePath)

		// Generate prompt
		buildPrompt(templatePath, configPath, outputPath)
	},
}

// Init function to add promptCmd to root
func init() {
	// Define CLI flags
	promptCmd.Flags().StringVarP(&topic, "topic", "", "", "Topic name (required)")
	promptCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to template file (default: resources/template.yaml)")
	promptCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file (default: based on topic)")
	promptCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to save generated prompt (default: based on topic)")

	promptCmd.MarkFlagRequired("topic")

}

// Builds the prompt and returns the path to the generated prompt
func buildPrompt(templatePath, configPath, outputPath string) string {
	// Generate prompt
	prompt.Build(templatePath, configPath, outputPath)
	return outputPath
}
