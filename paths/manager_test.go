package paths

import (
	"fmt"
	"testing"
)

func TestGetConfigPath_Custom(t *testing.T) {
	custom := "/custom/config.yaml"
	result := GetConfigPath("anytopic", custom)
	if result != custom {
		t.Errorf("Expected custom path %q, got %q", custom, result)
	}
}

func TestGetConfigPath_Default(t *testing.T) {
	topic := "testtopic"
	expected := fmt.Sprintf(ConfigPathFormat, topic)
	result := GetConfigPath(topic, "")
	if result != expected {
		t.Errorf("Expected default path %q, got %q", expected, result)
	}
}

func TestGetOutputPath_Custom(t *testing.T) {
	custom := "/custom/output.txt"
	result := GetOutputPath("anytopic", custom)
	if result != custom {
		t.Errorf("Expected custom path %q, got %q", custom, result)
	}
}

func TestGetOutputPath_Default(t *testing.T) {
	topic := "testtopic"
	expected := fmt.Sprintf(OutputPathFormat, topic)
	result := GetOutputPath(topic, "")
	if result != expected {
		t.Errorf("Expected default path %q, got %q", expected, result)
	}
}

func TestGetAnswerPath_Custom(t *testing.T) {
	custom := "/custom/answer.md"
	result := GetAnswerPath("anytopic", custom)
	if result != custom {
		t.Errorf("Expected custom path %q, got %q", custom, result)
	}
}

func TestGetAnswerPath_Default(t *testing.T) {
	topic := "testtopic"
	expected := fmt.Sprintf(AnswerPathFormat, topic)
	result := GetAnswerPath(topic, "")
	if result != expected {
		t.Errorf("Expected default path %q, got %q", expected, result)
	}
}

func TestGetTemplatePath_Custom(t *testing.T) {
	custom := "/custom/topic.yaml"
	result := GetTemplatePath(custom)
	if result != custom {
		t.Errorf("Expected custom template path %q, got %q", custom, result)
	}
}

func TestGetTemplatePath_Default(t *testing.T) {
	expected := TemplateFilePath
	result := GetTemplatePath("")
	if result != expected {
		t.Errorf("Expected default template path %q, got %q", expected, result)
	}
}
