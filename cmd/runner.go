package cmd

import (
	"fmt"
	"io"
	"log"
)

type ChatRunner struct {
	Out               io.Writer
	BuildPrompt       func(templatePath, configPath, outputPath string) string
	GetPrompt         func(promptFile string) (string, error)
	RunLLMInteraction func(prompt string) (string, error)
	SaveResponse      func(response, filePath string) error // ✅ corrected signature
}

func (cr *ChatRunner) Run() {
	fmt.Fprintln(cr.Out, "Generating prompt...")
	promptFilePath := cr.BuildPrompt(templatePath, configPath, outputPath)

	fmt.Fprintf(cr.Out, "Generated prompt saved to: %s\n", promptFilePath)

	fmt.Fprintln(cr.Out, "Reading generated prompt...")
	promptText, err := cr.GetPrompt(promptFilePath)
	if err != nil {
		log.Fatalf("Error retrieving prompt: %v", err)
	}

	fmt.Fprintln(cr.Out, "Generating LLM response...")
	response, err := cr.RunLLMInteraction(promptText)
	if err != nil {
		log.Fatalf("Error generating LLM response: %v", err)
	}

	fmt.Fprintf(cr.Out, "\nLLM Response:\n%s\n", response)

	if topic == "" {
		return
	}

	fmt.Fprintf(cr.Out, "Saving response to: %s\n", responseFilePath)
	if err := cr.SaveResponse(response, responseFilePath); err != nil {
		log.Fatalf("Error saving response: %v", err)
	}
}
