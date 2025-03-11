package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the main CLI command
var rootCmd = &cobra.Command{
	Use:   "topic-explorer",
	Short: "Generate a structured prompt from YAML templates",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Init function to register all subcommands
func init() {
	// Add all subcommands here
	rootCmd.AddCommand(llmCmd)    // Registers LLM command
	rootCmd.AddCommand(promptCmd) // Registers Prompt command
}
