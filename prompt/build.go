package prompt

import (
	"log"
	"os"

	"github.com/flosch/pongo2/v6"
	promptConfig "raja.aiml/ai.explorer/config/prompt"
	"raja.aiml/ai.explorer/paths"
)

func BuildTopicPrompt(templateFile, configFile, outputFile string) {
	log.Printf("[topic] Loading template: %s", templateFile)
	tpl := mustReadTemplate(templateFile)

	log.Printf("[topic] Loading config: %s", configFile)
	cfg := mustReadTopicConfig(configFile)

	ctx := pongo2.Context{
		"audience":                 cfg.Audience,
		"learning_stage":           cfg.LearningStage,
		"topic":                    cfg.Topic,
		"context":                  cfg.Context,
		"analogies":                cfg.Analogies,
		"concepts":                 cfg.Concepts,
		"explanation_requirements": cfg.ExplanationRequirements,
		"formatting":               cfg.Formatting,
		"constraints":              cfg.Constraints,
		"output_format":            cfg.OutputFormat,
		"purpose":                  cfg.Purpose,
		"tone":                     cfg.Tone,
	}

	renderAndSave(tpl.Template, ctx, outputFile)
}

// -------------------- Internal Helpers --------------------

func mustReadTemplate(path string) promptConfig.Template {
	tpl, err := promptConfig.ReadTemplate(path)
	if err != nil {
		log.Fatalf("Error reading template: %v", err)
	}
	return tpl
}

func mustReadTopicConfig(path string) promptConfig.TopicConfig {
	cfg, err := promptConfig.ReadTopicConfig(path)
	if err != nil {
		log.Fatalf("Error reading topic config: %v", err)
	}
	return cfg
}

func renderAndSave(tplStr string, ctx pongo2.Context, outputPath string) {
	log.Println("[render] Parsing template...")
	tpl, err := pongo2.FromString(tplStr)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	log.Println("[render] Executing template...")
	output, err := tpl.Execute(ctx)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	writePrompt(outputPath, output)
}

func writePrompt(path, content string) {
	log.Printf("[output] Writing to: %s", path)
	paths.EnsureDirectoryExists(path)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
	log.Printf("[output] Prompt generated successfully: %s", path)
}
