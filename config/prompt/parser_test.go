package prompt

import (
	"errors"
	"strings"
	"testing"
)

const (
	testTemplateYAML = `template: 'Hello, test!'`
	testInvalidYAML  = `invalid: [unclosed`
	testTopicYAML    = `
audience: "Developers"
learning_stage: "Intermediate"
topic: "Testing"
context: "Go"
analogies: "as a dependency injection example"
concepts: ["DI", "Testing"]
explanation_requirements: ["Explain DI"]
formatting: ["bullet points"]
constraints: ["none"]
output_format: ["markdown"]
purpose: "demonstrate generics"
tone: "informal"
`
)

func overrideReadFile(t *testing.T, mock func(string) ([]byte, error)) {
	t.Helper()
	original := readFile
	readFile = mock
	t.Cleanup(func() { readFile = original })
}

func TestReadYAML_Success(t *testing.T) {
	overrideReadFile(t, func(string) ([]byte, error) {
		return []byte(testTemplateYAML), nil
	})

	result, err := ReadYAML[Template]("dummy.yaml")
	assertNoError(t, err)
	assertEqual(t, result.Template, "Hello, test!", "template text")
}

func TestReadYAML_FileNotFound(t *testing.T) {
	expectedErr := errors.New("file not found")

	overrideReadFile(t, func(string) ([]byte, error) {
		return nil, expectedErr
	})

	_, err := ReadYAML[Template]("notfound.yaml")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %q, got %v", expectedErr.Error(), err)
	}
}

func TestReadYAML_UnmarshalError(t *testing.T) {
	overrideReadFile(t, func(string) ([]byte, error) {
		return []byte(testInvalidYAML), nil
	})

	_, err := ReadYAML[Template]("bad.yaml")
	if err == nil || !strings.Contains(err.Error(), "did not find expected") {
		t.Errorf("Expected unmarshal error, got: %v", err)
	}
}

func TestReadTemplate(t *testing.T) {
	overrideReadFile(t, func(string) ([]byte, error) {
		return []byte("template: 'Test Template'"), nil
	})

	result, err := ReadTemplate("dummy.yaml")
	assertNoError(t, err)
	assertEqual(t, result.Template, "Test Template", "template")
}

func TestReadTopicConfig(t *testing.T) {
	overrideReadFile(t, func(string) ([]byte, error) {
		return []byte(testTopicYAML), nil
	})

	cfg, err := ReadTopicConfig("dummy.yaml")
	assertNoError(t, err)

	assertEqual(t, cfg.Audience, "Developers", "Audience")
	assertEqual(t, cfg.Topic, "Testing", "Topic")
	assertEqual(t, len(cfg.Concepts), 2, "Concept count")
	assertEqual(t, cfg.Concepts[0], "DI", "First concept")
}

// -------------------- Assertion Helpers --------------------

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func assertEqual[T comparable](t *testing.T, got, want T, label string) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %s to be %v, got %v", label, want, got)
	}
}
