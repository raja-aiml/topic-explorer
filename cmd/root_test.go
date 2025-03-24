package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// --- Supportable for testable os.Exit() mocking ---
var osExit = os.Exit

func TestExecuteSuccess(t *testing.T) {
	oldArgs := os.Args
	os.Args = []string{"topic-explorer"}
	defer func() { os.Args = oldArgs }()

	rootCmd = &cobra.Command{
		Use: "topic-explorer",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "CLI ran successfully")
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Capture output using a buffer
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	Execute()

	assert.Contains(t, buf.String(), "CLI ran successfully")
}
