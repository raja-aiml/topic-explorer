package prompt

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Package-level variable for dependency injection.
var readFile = os.ReadFile

// Template structure for YAML.
type Template struct {
	Template string `yaml:"template"`
}

// Config structure for YAML.
type Config struct {
	Audience                string   `yaml:"audience"`
	LearningStage           string   `yaml:"learning_stage"`
	Topic                   string   `yaml:"topic"`
	Context                 string   `yaml:"context"`
	Analogies               string   `yaml:"analogies"`
	Concepts                []string `yaml:"concepts"`
	ExplanationRequirements []string `yaml:"explanation_requirements"`
	Formatting              []string `yaml:"formatting"`
	Constraints             []string `yaml:"constraints"`
	OutputFormat            []string `yaml:"output_format"`
	Purpose                 string   `yaml:"purpose"`
	Tone                    string   `yaml:"tone"`
}

// ReadYAML is a generic function that reads a YAML file into any type.
func ReadYAML[T any](filePath string) (T, error) {
	var result T
	data, err := readFile(filePath)
	if err != nil {
		return result, err
	}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil
}

// ReadTemplate parses the template YAML file using the generic parser.
func ReadTemplate(filePath string) (Template, error) {
	return ReadYAML[Template](filePath)
}

// ReadConfig parses the config YAML file using the generic parser.
func ReadConfig(filePath string) (Config, error) {
	return ReadYAML[Config](filePath)
}
