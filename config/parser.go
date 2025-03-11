package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Template structure for YAML
type Template struct {
	Template string `yaml:"template"`
}

// Config structure for YAML
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

// ReadTemplate parses the template YAML file
func ReadTemplate(filePath string) (Template, error) {
	var template Template
	data, err := os.ReadFile(filePath)
	if err != nil {
		return template, err
	}

	if err := yaml.Unmarshal(data, &template); err != nil {
		return template, err
	}

	return template, nil
}

// ReadConfig parses the config YAML file
func ReadConfig(filePath string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
