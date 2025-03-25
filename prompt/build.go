package prompt

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	promptConfig "raja.aiml/ai.explorer/config/prompt"
	"raja.aiml/ai.explorer/paths"
)

// Build orchestrates reading, processing, and writing prompt files
func Build(templateFile, configFile, outputFile string) {
	done := make(chan bool)
	go buildAsync(templateFile, configFile, outputFile, done)
	<-done
}

func buildAsync(templateFile, configFile, outputFile string, done chan<- bool) {
	var wg sync.WaitGroup

	templateCh := make(chan promptConfig.Template, 1)
	configCh := make(chan promptConfig.Config, 1)
	outputCh := make(chan string, 1)

	wg.Add(1)
	go asyncTask(func() {
		template, configData := loadYAMLFiles(templateFile, configFile)
		templateCh <- template
		configCh <- configData
		close(templateCh)
		close(configCh)
	}, &wg)

	template := <-templateCh
	configData := <-configCh

	wg.Add(1)
	go asyncTask(func() {
		outputCh <- generatePrompt(template.Template, configData)
		close(outputCh)
	}, &wg)

	promptOutput := <-outputCh

	wg.Add(1)
	go asyncTask(func() {
		savePrompt(outputFile, promptOutput)
	}, &wg)

	wg.Wait()
	done <- true
}

func asyncTask(task func(), wg *sync.WaitGroup) {
	defer wg.Done()
	task()
}

func loadYAMLFiles(templateFile, configFile string) (promptConfig.Template, promptConfig.Config) {
	template, err := promptConfig.ReadTemplate(templateFile)
	if err != nil {
		log.Fatalf("Error reading template: %v", err)
	}

	configData, err := promptConfig.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	return template, configData
}

func generatePrompt(templateStr string, configData promptConfig.Config) string {
	replacements := getReplacements(configData)

	output := templateStr
	for key, value := range replacements {
		output = strings.ReplaceAll(output, key, value)
	}

	return output
}

func getReplacements(configData promptConfig.Config) map[string]string {
	replacements := make(map[string]string)
	val := reflect.ValueOf(configData)
	typ := reflect.TypeOf(configData)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("yaml")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		} else {
			tag = strings.Split(tag, ",")[0]
		}

		placeholder := "{" + tag + "}"
		fieldValue := val.Field(i)
		var strValue string

		if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.String {
			var items []string
			for j := 0; j < fieldValue.Len(); j++ {
				items = append(items, fieldValue.Index(j).String())
			}
			strValue = paths.FormatList(items)
		} else if fieldValue.Kind() == reflect.String {
			strValue = fieldValue.String()
		} else {
			strValue = fmt.Sprintf("%v", fieldValue.Interface())
		}

		replacements[placeholder] = strValue
	}

	return replacements
}

func savePrompt(outputFile, promptOutput string) {
	paths.EnsureDirectoryExists(outputFile)

	if err := os.WriteFile(outputFile, []byte(promptOutput), 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Prompt generated successfully:", outputFile)
}
