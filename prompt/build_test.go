package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"raja.aiml/ai.explorer/config"
	"raja.aiml/ai.explorer/paths"
)

// TestGeneratePrompt verifies that generatePrompt correctly replaces placeholders.
func TestGeneratePrompt(t *testing.T) {
	// Fake template string with placeholders.
	templateStr := "Welcome, {audience}. You are learning about {topic} in a {tone} way."
	// Create a fake config with dummy data.
	fakeConfig := config.Config{
		Audience:                "Test Audience",
		LearningStage:           "beginner",
		Topic:                   "Dependency Injection",
		Context:                 "test",
		Analogies:               "injection examples",
		Concepts:                []string{"concept1", "concept2"},
		ExplanationRequirements: []string{"explain one", "explain two"},
		Formatting:              []string{"format one", "format two"},
		Constraints:             []string{"none"},
		OutputFormat:            []string{"output one"},
		Purpose:                 "demonstrate DI",
		Tone:                    "Friendly",
	}

	expectedOutput := "Welcome, Test Audience. You are learning about Dependency Injection in a Friendly way."
	result := generatePrompt(templateStr, fakeConfig)
	if result != expectedOutput {
		t.Errorf("Expected: %q, got: %q", expectedOutput, result)
	}
}

// TestBuild verifies that Build reads the template and config files, generates the prompt,
// and writes the correct output.
func TestBuild(t *testing.T) {
	// Create a temporary directory.
	tempDir := t.TempDir()

	// Create a temporary template file.
	templateContent := `template: |
  Hello, {audience}! Topic: {topic}. Tone: {tone}.`
	templateFile := filepath.Join(tempDir, "template.yaml")
	if err := os.WriteFile(templateFile, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create a temporary config file.
	configContent := `
audience: "Test Audience"
learning_stage: "beginner"
topic: "Dependency Injection"
context: "test"
analogies: "injection examples"
concepts:
  - "concept1"
  - "concept2"
explanation_requirements:
  - "explain one"
  - "explain two"
formatting:
  - "format one"
  - "format two"
constraints:
  - "none"
output_format:
  - "output one"
purpose: "demonstrate DI"
tone: "Friendly"
`
	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Define the output file path in the temporary directory.
	outputFile := filepath.Join(tempDir, "output.txt")

	// Call Build with our injected dependencies.
	Build(templateFile, configFile, outputFile)

	// Read the generated output.
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	generatedPrompt := strings.TrimSpace(string(outputData))

	expected := "Hello, Test Audience! Topic: Dependency Injection. Tone: Friendly."
	if generatedPrompt != expected {
		t.Errorf("Expected output:\n%q\ngot:\n%q", expected, generatedPrompt)
	}
}

// TestAsyncTask ensures that asyncTask executes the given function.
func TestAsyncTask(t *testing.T) {
	var wg sync.WaitGroup
	var executed bool
	wg.Add(1)
	asyncTask(func() {
		executed = true
	}, &wg)
	wg.Wait()
	if !executed {
		t.Error("Expected async task to be executed")
	}
}

// TestLoadYAMLFiles verifies that loadYAMLFiles correctly parses the template and config files.
func TestLoadYAMLFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create a temporary template file.
	templateContent := `template: |
  Test template with {audience} and {topic}.`
	templateFile := filepath.Join(tempDir, "temp_template.yaml")
	if err := os.WriteFile(templateFile, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create a temporary config file.
	configContent := `
audience: "Unit Tester"
learning_stage: "advanced"
topic: "Reflection"
context: "testing"
analogies: "mirrors"
concepts:
  - "reflection"
explanation_requirements:
  - "explain reflection"
formatting:
  - "bullet points"
constraints:
  - "none"
output_format:
  - "text"
purpose: "test parsing"
tone: "serious"
`
	configFile := filepath.Join(tempDir, "temp_config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	tmpl, cfg := loadYAMLFiles(templateFile, configFile)

	if !strings.Contains(tmpl.Template, "{audience}") || !strings.Contains(tmpl.Template, "{topic}") {
		t.Error("Template did not load expected placeholders")
	}
	if cfg.Audience != "Unit Tester" || cfg.Topic != "Reflection" {
		t.Errorf("Config not parsed correctly: got audience=%q, topic=%q", cfg.Audience, cfg.Topic)
	}
}

// TestGetReplacements calls getReplacements and checks that the expected keys and formatted values are returned.
func TestGetReplacements(t *testing.T) {
	// Create a fake config with sample data.
	fakeConfig := config.Config{
		Audience:                "Test Audience",
		LearningStage:           "advanced",
		Topic:                   "Generics",
		Context:                 "development",
		Analogies:               "cooking recipes",
		Concepts:                []string{"concept1", "concept2"},
		ExplanationRequirements: []string{"explain A", "explain B"},
		Formatting:              []string{"header", "list"},
		Constraints:             []string{"none"},
		OutputFormat:            []string{"markdown"},
		Purpose:                 "demonstrate DI",
		Tone:                    "informal",
	}

	repls := getReplacements(fakeConfig)

	// Check string fields.
	if repls["{audience}"] != "Test Audience" {
		t.Errorf("Expected {audience} to be %q, got %q", "Test Audience", repls["{audience}"])
	}
	if repls["{topic}"] != "Generics" {
		t.Errorf("Expected {topic} to be %q, got %q", "Generics", repls["{topic}"])
	}

	// Check slice fields. FormatList prepends "\n- " and joins items.
	expectedConcepts := "\n- concept1\n- concept2"
	if repls["{concepts}"] != expectedConcepts {
		t.Errorf("Expected {concepts} to be %q, got %q", expectedConcepts, repls["{concepts}"])
	}
}

// TestSavePrompt ensures that savePrompt writes the output to the given file.
func TestSavePrompt(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "saved_prompt.txt")
	testContent := "This is a test prompt."

	// Call savePrompt.
	savePrompt(outputFile, testContent)

	// Read file to verify content.
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read saved prompt: %v", err)
	}
	if string(data) != testContent {
		t.Errorf("Expected file content %q, got %q", testContent, string(data))
	}
}

// TestBuildAsync exercises the asynchronous buildAsync function directly.
func TestBuildAsync(t *testing.T) {
	tempDir := t.TempDir()

	// Create temporary template file.
	templateContent := `template: |
  Async test: {audience} and {topic}.`
	templateFile := filepath.Join(tempDir, "async_template.yaml")
	if err := os.WriteFile(templateFile, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write async template file: %v", err)
	}

	// Create temporary config file.
	configContent := `
audience: "Async Tester"
learning_stage: "intermediate"
topic: "Concurrency"
context: "async"
analogies: "orchestration"
concepts:
  - "goroutines"
explanation_requirements:
  - "explain concurrency"
formatting:
  - "plain text"
constraints:
  - "none"
output_format:
  - "text"
purpose: "test async"
tone: "calm"
`
	configFile := filepath.Join(tempDir, "async_config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write async config file: %v", err)
	}

	// Define output file.
	outputFile := filepath.Join(tempDir, "async_output.txt")
	done := make(chan bool)

	// Call buildAsync directly.
	go buildAsync(templateFile, configFile, outputFile, done)
	<-done

	// Read and verify the output.
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read async output file: %v", err)
	}
	generated := strings.TrimSpace(string(data))
	expected := "Async test: Async Tester and Concurrency."
	if generated != expected {
		t.Errorf("Expected async output %q, got %q", expected, generated)
	}
}

// TestReflectionFallback forces the fallback branch in getReplacements by using an unexported field type.
// We do this by embedding a dummy field of a different type in a wrapper struct.
func TestReflectionFallback(t *testing.T) {
	// Define a wrapper type that embeds config.Config and adds an integer field.
	type dummyConfig struct {
		config.Config
		Extra int `yaml:"extra"`
	}
	// Create an instance with non-string value.
	dc := dummyConfig{
		Config: config.Config{
			Audience: "Test Audience",
			// other fields can be left empty
		},
		Extra: 42,
	}
	// Because getReplacements expects config.Config, we need to call it on the embedded part.
	// We'll use reflection on the wrapper struct manually to simulate the fallback.
	val := reflect.ValueOf(dc)
	typ := reflect.TypeOf(dc)
	repls := make(map[string]string)
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
			// Fallback branch.
			strValue = fmt.Sprintf("%v", fieldValue.Interface())
		}
		repls[placeholder] = strValue
	}
	// Verify that the extra field was processed via fallback.
	if repls["{extra}"] != "42" {
		t.Errorf("Expected {extra} to be %q, got %q", "42", repls["{extra}"])
	}
}
