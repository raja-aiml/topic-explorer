package prompt

import (
	"fmt"
	"log"
	"os"

	"github.com/flosch/pongo2/v6"
	promptConfig "raja.aiml/ai.explorer/config/prompt"
	"raja.aiml/ai.explorer/paths"
)

func BuildTopicPrompt(templateFile, configFile, outputFile string) {
	template := mustReadTemplate(templateFile)
	config := mustReadTopicConfig(configFile)
	ctx := mustBuildTopicContext(config)
	renderAndSave(template.Template, ctx, outputFile)
}

func BuildChartPrompt(templateFile, configFile, outputFile string) {
	template := mustReadTemplate(templateFile)
	config := mustReadChartConfig(configFile)
	ctx := mustBuildChartContext(config)
	renderAndSave(template.Template, ctx, outputFile)
}

// -------------------- Internal Helpers --------------------

func mustReadTemplate(file string) promptConfig.Template {
	tpl, err := promptConfig.ReadTemplate(file)
	if err != nil {
		log.Fatalf("Error reading template: %v", err)
	}
	return tpl
}

func mustReadTopicConfig(file string) promptConfig.TopicConfig {
	cfg, err := promptConfig.ReadTopicConfig(file)
	if err != nil {
		log.Fatalf("Error reading topic config: %v", err)
	}
	return cfg
}

func mustReadChartConfig(file string) promptConfig.ChartConfig {
	cfg, err := promptConfig.ReadChartConfig(file)
	if err != nil {
		log.Fatalf("Error reading chart config: %v", err)
	}
	return cfg
}

func mustBuildTopicContext(config promptConfig.TopicConfig) pongo2.Context {
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
	}
}

func mustBuildChartContext(config promptConfig.ChartConfig) pongo2.Context {
	return pongo2.Context{
		"flow_direction":  config.FlowDirection,
		"planning_phase":  config.PlanningPhase,
		"execution_phase": config.ExecutionPhase,
		"planning_links":  config.GenerateLinks(config.PlanningPhase.Steps),
		"execution_links": config.GenerateLinks(config.ExecutionPhase.Steps),
		"transition_link": map[string]string{
			"from": config.PlanningPhase.Steps[len(config.PlanningPhase.Steps)-1].ID,
			"to":   config.ExecutionPhase.Steps[0].ID,
		},
		"style": config.Style,
	}
}

func renderAndSave(tplStr string, ctx pongo2.Context, outputFile string) {
	tpl, err := pongo2.FromString(tplStr)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	output, err := tpl.Execute(ctx)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}
	savePrompt(outputFile, output)
}

func savePrompt(path, content string) {
	paths.EnsureDirectoryExists(path)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("Error writing output: %v", err)
	}
	fmt.Println("Prompt generated successfully:", path)
}
