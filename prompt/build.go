package prompt

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"raja.aiml/ai.explorer/config"
	"raja.aiml/ai.explorer/paths"
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

// getReplacements uses reflection to create a map of placeholder replacements
// based on the YAML tags in the config.Config struct.
func getReplacements(configData config.Config) map[string]string {
	replacements := make(map[string]string)
	val := reflect.ValueOf(configData)
	typ := reflect.TypeOf(configData)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Use the YAML tag if available; otherwise, default to the lowercased field name.
		tag := field.Tag.Get("yaml")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		} else {
			// In case the tag has additional options (e.g. "audience,omitempty"),
			// split and take the first value.
			tag = strings.Split(tag, ",")[0]
		}

		placeholder := "{" + tag + "}"
		fieldValue := val.Field(i)
		var strValue string

		// If the field is a slice of strings, use paths.FormatList to format it.
		if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.String {
			var items []string
			for j := 0; j < fieldValue.Len(); j++ {
				items = append(items, fieldValue.Index(j).String())
			}
			strValue = paths.FormatList(items)
		} else if fieldValue.Kind() == reflect.String {
			strValue = fieldValue.String()
		} else {
			// Fallback for other types: use the default string formatting.
			strValue = fmt.Sprintf("%v", fieldValue.Interface())
		}

		replacements[placeholder] = strValue
	}

	return replacements
}

// savePrompt ensures directories exist and writes the output file
func savePrompt(outputFile, promptOutput string) {
	paths.EnsureDirectoryExists(outputFile)

	if err := os.WriteFile(outputFile, []byte(promptOutput), 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Prompt generated successfully:", outputFile)
}
