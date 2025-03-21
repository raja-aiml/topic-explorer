package config

import (
	"errors"
	"strings"
	"testing"
)

func TestReadYAMLSuccess(t *testing.T) {
	// Save the original readFile and restore it after the test.
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	// Create dummy YAML content for a Template.
	yamlContent := []byte("template: 'Hello, test!'")
	readFile = func(filename string) ([]byte, error) {
		return yamlContent, nil
	}

	result, err := ReadYAML[Template]("dummy.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Template != "Hello, test!" {
		t.Errorf("Expected template 'Hello, test!', got %q", result.Template)
	}
}

func TestReadYAMLFileNotFound(t *testing.T) {
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	// Override readFile to return an error.
	expectedErr := errors.New("file not found")
	readFile = func(filename string) ([]byte, error) {
		return nil, expectedErr
	}

	_, err := ReadYAML[Template]("nonexistent.yaml")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %q, got %q", expectedErr.Error(), err.Error())
	}
}

func TestReadYAMLUnmarshalError(t *testing.T) {
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	// Override readFile to return invalid YAML.
	readFile = func(filename string) ([]byte, error) {
		return []byte("invalid: [unclosed"), nil
	}

	_, err := ReadYAML[Template]("dummy.yaml")
	if err == nil {
		t.Fatal("Expected unmarshal error, got nil")
	}
	// Check that the error message contains an expected substring.
	if !strings.Contains(err.Error(), "did not find expected") {
		t.Errorf("Unexpected unmarshal error: %v", err)
	}
}

func TestReadTemplate(t *testing.T) {
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	yamlContent := []byte("template: 'Test Template'")
	readFile = func(filename string) ([]byte, error) {
		return yamlContent, nil
	}

	result, err := ReadTemplate("dummy.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Template != "Test Template" {
		t.Errorf("Expected 'Test Template', got %q", result.Template)
	}
}

func TestReadConfig(t *testing.T) {
	origReadFile := readFile
	defer func() { readFile = origReadFile }()

	// Create a dummy YAML config.
	yamlContent := []byte(`
audience: "Developers"
learning_stage: "Intermediate"
topic: "Testing"
context: "Go"
analogies: "as a dependency injection example"
concepts:
  - "DI"
  - "Testing"
explanation_requirements:
  - "Explain DI"
formatting:
  - "bullet points"
constraints:
  - "none"
output_format:
  - "markdown"
purpose: "demonstrate generics"
tone: "informal"
`)
	readFile = func(filename string) ([]byte, error) {
		return yamlContent, nil
	}

	result, err := ReadConfig("dummy.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Audience != "Developers" {
		t.Errorf("Expected Audience 'Developers', got %q", result.Audience)
	}
	if result.Topic != "Testing" {
		t.Errorf("Expected Topic 'Testing', got %q", result.Topic)
	}
	if len(result.Concepts) != 2 {
		t.Errorf("Expected 2 concepts, got %d", len(result.Concepts))
	}
	if result.Concepts[0] != "DI" {
		t.Errorf("Expected first concept 'DI', got %q", result.Concepts[0])
	}
}
