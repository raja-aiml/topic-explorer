package paths

import "fmt"

// Default paths
const (
	BasePath         = "resources/templates"
	ConfigPathFormat = BasePath + "/configs/%s.yaml"
	OutputPathFormat = BasePath + "/output/%s/prompt.txt"
	AnswerPathFormat = BasePath + "/output/%s/answer.md"
	TemplateFilePath = BasePath + "/topic.yaml"
)

// GetConfigPath returns the config file path for a given topic
func GetConfigPath(topic string, customPath string) string {
	if customPath != "" {
		return customPath
	}
	return fmt.Sprintf(ConfigPathFormat, topic)
}

// GetOutputPath returns the output file path for a given topic
func GetOutputPath(topic string, customPath string) string {
	if customPath != "" {
		return customPath
	}
	return fmt.Sprintf(OutputPathFormat, topic)
}

// GetAnswerPath returns the answer file path for a given topic
func GetAnswerPath(topic string, customPath string) string {
	if customPath != "" {
		return customPath
	}
	return fmt.Sprintf(AnswerPathFormat, topic)
}

// GetTemplatePath returns the template file path (default unless overridden)
func GetTemplatePath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	return TemplateFilePath
}
