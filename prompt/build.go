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

	log.Printf("[topic] Rendering context: %+v", cfg)

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

func BuildChartPrompt(templateFile, configFile, outputFile string) {
	log.Printf("[chart] Loading template: %s", templateFile)
	tpl := mustReadTemplate(templateFile)

	log.Printf("[chart] Loading config: %s", configFile)
	cfg := mustReadChartConfig(configFile)

	log.Println("[chart] Generating links...")
	cfg.PlanningLinks = cfg.GenerateLinks(cfg.PlanningPhase.Steps)
	cfg.ExecutionLinks = cfg.GenerateLinks(cfg.ExecutionPhase.Steps)

	if len(cfg.PlanningPhase.Steps) > 0 && len(cfg.ExecutionPhase.Steps) > 0 {
		cfg.TransitionLink = promptConfig.Transition{
			From: cfg.PlanningPhase.Steps[len(cfg.PlanningPhase.Steps)-1].ID,
			To:   cfg.ExecutionPhase.Steps[0].ID,
		}
	}

	log.Printf("[chart] Rendering context:\n  Flow: %s\n  Planning Steps: %d\n  Execution Steps: %d\n  Transition: %s -> %s\n",
		cfg.FlowDirection,
		len(cfg.PlanningPhase.Steps),
		len(cfg.ExecutionPhase.Steps),
		cfg.TransitionLink.From,
		cfg.TransitionLink.To,
	)

	ctx := pongo2.Context{
		"flow_direction":  cfg.FlowDirection,
		"planning_phase":  cfg.PlanningPhase,
		"execution_phase": cfg.ExecutionPhase,
		"planning_links":  cfg.PlanningLinks,
		"execution_links": cfg.ExecutionLinks,
		"transition_link": cfg.TransitionLink,
		"style":           cfg.Style,
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

func mustReadChartConfig(path string) promptConfig.ChartConfig {
	cfg, err := promptConfig.ReadChartConfig(path)
	if err != nil {
		log.Fatalf("Error reading chart config: %v", err)
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

	log.Println(tpl)
	log.Println(output)

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
