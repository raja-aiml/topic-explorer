package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	topicTemplateYAML = `template: |
  Hello, {{ audience }}!
  You are learning about {{ topic }} in a {{ tone }} way.`

	topicConfigYAML = `
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

	expectedOutput = `Hello, Test Audience!
You are learning about Generics in a friendly way.`
)

func TestBuildTopicPrompt(t *testing.T) {
	dir := t.TempDir()

	tplPath := filepath.Join(dir, "template.yaml")
	cfgPath := filepath.Join(dir, "config.yaml")
	outPath := filepath.Join(dir, "output.txt")

	writeFile(t, tplPath, topicTemplateYAML)
	writeFile(t, cfgPath, topicConfigYAML)

	BuildTopicPrompt(tplPath, cfgPath, outPath)

	got := readFile(t, outPath)

	if strings.TrimSpace(got) != expectedOutput {
		t.Errorf("\nExpected:\n%q\nGot:\n%q", expectedOutput, got)
	}
}

// Helpers

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile failed: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile failed: %v", err)
	}
	return string(data)
}
