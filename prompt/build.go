package prompt

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"raja.aiml/topic.explorer/config"
	"raja.aiml/topic.explorer/paths"
)

// ProcessFiles orchestrates reading, processing, and writing prompt files
func Build(templateFile, configFile, outputFile string) {
	// Create a channel to receive completion signal
	done := make(chan bool)
	// Generate prompt asynchronously
	go buildAsync(templateFile, configFile, outputFile, done)
	// Wait for completion signal
	<-done
}

// BuildAsync orchestrates reading, processing, and writing prompt files asynchronously
func buildAsync(templateFile, configFile, outputFile string, done chan<- bool) {
	var wg sync.WaitGroup

	// Channels for async operations
	templateCh := make(chan config.Template, 1)
	configCh := make(chan config.Config, 1)
	outputCh := make(chan string, 1)

	// Start async YAML file reading
	wg.Add(1)
	go asyncTask(func() {
		template, configData := loadYAMLFiles(templateFile, configFile)
		templateCh <- template
		configCh <- configData
		close(templateCh)
		close(configCh)
	}, &wg)

	// Read results from channels
	template := <-templateCh
	configData := <-configCh

	// Start async prompt generation
	wg.Add(1)
	go asyncTask(func() {
		outputCh <- generatePrompt(template.Template, configData)
		close(outputCh)
	}, &wg)

	// Read generated prompt output
	promptOutput := <-outputCh

	// Start async file saving
	wg.Add(1)
	go asyncTask(func() {
		savePrompt(outputFile, promptOutput)
	}, &wg)

	// Wait for all operations to complete
	wg.Wait()

	// Signal completion
	done <- true
}

// asyncTask runs a function inside a goroutine and tracks it with WaitGroup
func asyncTask(task func(), wg *sync.WaitGroup) {
	defer wg.Done()
	task()
}

// loadYAMLFiles reads the YAML template and config
func loadYAMLFiles(templateFile, configFile string) (config.Template, config.Config) {
	template, err := config.ReadTemplate(templateFile)
	if err != nil {
		log.Fatalf("Error reading template: %v", err)
	}

	configData, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	return template, configData
}

// generatePrompt replaces placeholders in the template with actual config values
func generatePrompt(templateStr string, configData config.Config) string {
	replacements := getReplacements(configData)

	output := templateStr
	for key, value := range replacements {
		output = strings.ReplaceAll(output, key, value)
	}

	return output
}

// getReplacements returns a map of placeholder replacements
func getReplacements(configData config.Config) map[string]string {
	return map[string]string{
		"{audience}":                 configData.Audience,
		"{topic}":                    configData.Topic,
		"{learning_stage}":           configData.LearningStage,
		"{context}":                  configData.Context,
		"{analogies}":                configData.Analogies,
		"{concepts}":                 paths.FormatList(configData.Concepts),
		"{explanation_requirements}": paths.FormatList(configData.ExplanationRequirements),
		"{formatting}":               paths.FormatList(configData.Formatting),
		"{constraints}":              paths.FormatList(configData.Constraints),
		"{output_format}":            paths.FormatList(configData.OutputFormat),
		"{purpose}":                  configData.Purpose,
		"{tone}":                     configData.Tone,
	}
}

// savePrompt ensures directories exist and writes the output file
func savePrompt(outputFile, promptOutput string) {
	paths.EnsureDirectoryExists(outputFile)

	if err := os.WriteFile(outputFile, []byte(promptOutput), 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Prompt generated successfully:", outputFile)
}
