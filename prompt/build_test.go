package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raja.aiml/topic.explorer/config"
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

	// Expected output after replacements.
	expectedOutput := "Welcome, Test Audience. You are learning about Dependency Injection in a Friendly way."
	result := generatePrompt(templateStr, fakeConfig)
	if result != expectedOutput {
		t.Errorf("Expected: %q, got: %q", expectedOutput, result)
	}
}

// TestBuild verifies that Build reads the template and config files, generates the prompt,
// and writes the correct output. We inject temporary files to avoid hardcoded paths.
func TestBuild(t *testing.T) {
	// Create a temporary directory.
	tempDir := t.TempDir()

	// Create a temporary template file.
	// Note: the template YAML must include the "template:" key so that config.ReadTemplate unmarshals it.
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

	// Call Build with our injected dependencies (templateFile, configFile, and outputFile).
	Build(templateFile, configFile, outputFile)

	// Read the generated output.
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	generatedPrompt := strings.TrimSpace(string(outputData))

	// Expected output (matches the template with injected values).
	expected := "Hello, Test Audience! Topic: Dependency Injection. Tone: Friendly."
	if generatedPrompt != expected {
		t.Errorf("Expected output:\n%q\ngot:\n%q", expected, generatedPrompt)
	}
}
