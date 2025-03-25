package prompt

import (
	"fmt"
	"log"
	"os"

	"github.com/flosch/pongo2/v6"
	promptConfig "raja.aiml/ai.explorer/config/prompt"
	"raja.aiml/ai.explorer/paths"
)

// Build reads template and config, renders the prompt, and writes to file
func Build(templateFile, configFile, outputFile string) {
	// Step 1: Read YAML files
	template, err := promptConfig.ReadTemplate(templateFile)
	if err != nil {
		log.Fatalf("Error reading template: %v", err)
	}

	configData, err := promptConfig.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Step 2: Parse template with pongo2
	tpl, err := pongo2.FromString(template.Template)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Step 3: Convert config struct to pongo2 context
	ctx, err := toPongoContext(configData)
	if err != nil {
		log.Fatalf("Error converting config to context: %v", err)
	}

	// Step 4: Render
	output, err := tpl.Execute(ctx)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	// Step 5: Save result
	savePrompt(outputFile, output)
}

func toPongoContext(config promptConfig.Config) (pongo2.Context, error) {
	// We can manually convert or use reflection (here we assume a flat config)
	return pongo2.Context{
		"audience":                 config.Audience,
		"learning_stage":           config.LearningStage,
		"topic":                    config.Topic,
		"context":                  config.Context,
		"analogies":                config.Analogies,
		"concepts":                 config.Concepts,
		"explanation_requirements": config.ExplanationRequirements,
		"formatting":               config.Formatting,
		"constraints":              config.Constraints,
		"output_format":            config.OutputFormat,
		"purpose":                  config.Purpose,
		"tone":                     config.Tone,
	}, nil
}

func savePrompt(outputFile, promptOutput string) {
	paths.EnsureDirectoryExists(outputFile)

	if err := os.WriteFile(outputFile, []byte(promptOutput), 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Prompt generated successfully:", outputFile)
}
