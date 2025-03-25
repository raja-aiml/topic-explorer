package prompt

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Allows mocking in tests
var readFile = os.ReadFile

// Template holds a Pongo2-compatible string from YAML.
type Template struct {
	Template string `yaml:"template"`
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

// -------------------- Flowchart Prompt --------------------

type Step struct {
	ID          string `yaml:"id"`
	Emoji       string `yaml:"emoji"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type Phase struct {
	Title     string `yaml:"title"`
	Emoji     string `yaml:"emoji"`
	Direction string `yaml:"direction"`
	Steps     []Step `yaml:"steps"`
}

type Style struct {
	SubtaskBox string `yaml:"subtaskBox"`
}

type Link struct {
	From string
	To   string
}

type Transition struct {
	From string
	To   string
}

type ChartConfig struct {
	FlowDirection  string     `yaml:"flow_direction"`
	Style          Style      `yaml:"style"`
	PlanningPhase  Phase      `yaml:"planning_phase"`
	ExecutionPhase Phase      `yaml:"execution_phase"`
	PlanningLinks  []Link     `yaml:"-"`
	ExecutionLinks []Link     `yaml:"-"`
	TransitionLink Transition `yaml:"-"`
}

// -------------------- YAML Loaders --------------------

// ReadYAML loads any YAML file into the given type.
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

// ReadTemplate loads a prompt template from a YAML file.
func ReadTemplate(filePath string) (Template, error) {
	return ReadYAML[Template](filePath)
}

// ReadTopicConfig loads a topic-based config (e.g. Git prompt).
func ReadTopicConfig(filePath string) (TopicConfig, error) {
	return ReadYAML[TopicConfig](filePath)
}

// ReadChartConfig loads a flowchart-style config (e.g. Planning/Execution).
func ReadChartConfig(filePath string) (ChartConfig, error) {
	return ReadYAML[ChartConfig](filePath)
}

func (c ChartConfig) GenerateLinks(steps []Step) []Link {
	var links []Link
	for i := 0; i < len(steps)-1; i++ {
		links = append(links, Link{
			From: steps[i].ID,
			To:   steps[i+1].ID,
		})
	}
	return links
}
