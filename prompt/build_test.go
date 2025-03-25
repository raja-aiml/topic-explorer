package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTemplate = `template: |
  Hello, {{ audience }}!
  You are learning about {{ topic }} in a {{ tone }} way.`

const testConfig = `
audience: "Test Audience"
learning_stage: "beginner"
topic: "Generics"
context: "testing"
analogies: "boxes and types"
concepts: ["type parameters", "constraints"]
explanation_requirements: ["simple example"]
formatting: ["bullets"]
constraints: ["compile-time"]
output_format: ["markdown"]
purpose: "learn Go generics"
tone: "friendly"
`

var expectedText = `Hello, Test Audience!
You are learning about Generics in a friendly way.`

func TestBuildPromptCorrectly(t *testing.T) {
	dir := t.TempDir()

	templatePath := filepath.Join(dir, "template.yaml")
	configPath := filepath.Join(dir, "config.yaml")
	outputPath := filepath.Join(dir, "output.txt")

	writeFile(t, templatePath, testTemplate)
	writeFile(t, configPath, testConfig)

	Build(templatePath, configPath, outputPath)

	got := readFile(t, outputPath)
	want := expectedText

	if strings.TrimSpace(got) != want {
		t.Errorf("\nExpected:\n%q\nGot:\n%q", want, got)
	}
}

// writeFile writes string content to a file path.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile failed for %s: %v", path, err)
	}
}

// readFile reads file content as string.
func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile failed for %s: %v", path, err)
	}
	return string(data)
}
