package prompt

import (
	"os"

	"gopkg.in/yaml.v3"
)

// -------------------- File I/O --------------------

// Allows mocking in tests
var readFile = os.ReadFile

// Generic YAML loader
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

// -------------------- Template --------------------

type Template struct {
	Template string `yaml:"template"`
}

func ReadTemplate(filePath string) (Template, error) {
	return ReadYAML[Template](filePath)
}

// -------------------- Topic Prompt --------------------

type TopicConfig struct {
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

func ReadTopicConfig(filePath string) (TopicConfig, error) {
	return ReadYAML[TopicConfig](filePath)
}

// -------------------- Flowchart Prompt --------------------

type ChartConfig struct {
	FlowDirection  string            `yaml:"flow_direction"`
	Style          map[string]string `yaml:"style"`
	PlanningPhase  Phase             `yaml:"planning_phase"`
	PlanningLinks  []Link            `yaml:"planning_links"`
	ExecutionPhase Phase             `yaml:"execution_phase"`
	ExecutionLinks []Link            `yaml:"execution_links"`
	TransitionLink Link              `yaml:"transition_link"`
}

type Phase struct {
	Title     string `yaml:"title"`
	Emoji     string `yaml:"emoji"`
	Direction string `yaml:"direction"`
	Steps     []Step `yaml:"steps"`
}

type Step struct {
	ID          string `yaml:"id"`
	Emoji       string `yaml:"emoji"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type Link struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

func ReadChartConfig(filePath string) (ChartConfig, error) {
	return ReadYAML[ChartConfig](filePath)
}

// Optional utility for building links from steps (used in generators/test data)
func GenerateLinks(steps []Step) []Link {
	var links []Link
	for i := 0; i < len(steps)-1; i++ {
		links = append(links, Link{
			From: steps[i].ID,
			To:   steps[i+1].ID,
		})
	}
	return links
}
