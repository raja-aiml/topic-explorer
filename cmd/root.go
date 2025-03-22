package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the main CLI command
var rootCmd = &cobra.Command{
	Use:           "topic-explorer",
	Short:         "Generate a structured prompt from YAML templates and interact with LLMs",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the CLI application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
