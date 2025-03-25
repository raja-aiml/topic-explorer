package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "topic-explorer",
	Short:         "Prompt generation + LLM interaction CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	godotenv.Load()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(rootCmd.ErrOrStderr(), "Error: %v\n", err)
		os.Exit(1)
	}
}
