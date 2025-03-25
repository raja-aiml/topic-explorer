package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// --- Configurable Constants ---
const (
	binaryName  = "ai-explorer"
	topic       = "git"
	model       = "gpt-4o"
	provider    = "openai"
	temperature = "0.7"
)

// --- Paths Struct ---
type TestPaths struct {
	RootDir      string
	BinPath      string
	TemplatePath string
	ConfigPath   string
	OutputDir    string
	PromptOutput string
}

// --- Path Resolver ---
func newTestPaths(topic string) *TestPaths {
	root := ".."
	output := filepath.Join(root, ".build", "output", "test_"+topic)

	return &TestPaths{
		RootDir:      root,
		BinPath:      filepath.Join(".build", binaryName),
		TemplatePath: filepath.Join("resources", "template.yaml"),
		ConfigPath:   filepath.Join("resources", "configs", topic+".yaml"),
		OutputDir:    output,
		PromptOutput: filepath.Join(output, "prompt.txt"),
	}
}

// --- Command Executor ---
func runCommand(paths *TestPaths, args ...string) ([]byte, error) {
	cmd := exec.Command(paths.BinPath, args...)
	cmd.Dir = paths.RootDir // run from project root
	cmd.Env = os.Environ()  // inherit .env-loaded vars
	return cmd.CombinedOutput()
}

var _ = Describe("AI Explorer CLI (E2E)", func() {
	var paths *TestPaths

	BeforeEach(func() {
		paths = newTestPaths(topic)
		err := os.MkdirAll(paths.OutputDir, os.ModePerm)
		Expect(err).ToNot(HaveOccurred(), "Failed to create output directory")
	})

	Context("Given a valid YAML template and topic config", func() {
		When("the user runs the 'prompt' command", func() {
			It("Then it should generate a prompt file", func() {
				output, err := runCommand(paths,
					"prompt",
					"--topic", topic,
					"--template", paths.TemplatePath,
					"--config", paths.ConfigPath,
					"--output", paths.PromptOutput,
				)

				Expect(err).ToNot(HaveOccurred(), "Prompt generation failed:\n%s", string(output))
				Expect(paths.PromptOutput).To(BeAnExistingFile(), "Prompt file should exist")
			})
		})

		When("the user runs the 'llm' command with the generated prompt", func() {
			It("Then it should return a valid model response", func() {
				output, err := runCommand(paths,
					"llm",
					"--provider", provider,
					"--model", model,
					"--prompt", paths.PromptOutput,
					"--temperature", temperature,
				)

				Expect(err).ToNot(HaveOccurred(), "LLM command failed:\n%s", string(output))
				Expect(string(output)).To(ContainSubstring(topic))
			})
		})

		When("the user runs the 'chat' command", func() {
			It("Then it should generate a prompt and get a response in one step", func() {
				output, err := runCommand(paths,
					"chat",
					"--topic", topic,
					"--provider", provider,
					"--model", model,
				)

				Expect(err).ToNot(HaveOccurred(), "Chat command failed:\n%s", string(output))
				Expect(string(output)).To(ContainSubstring(topic))
			})
		})
	})
})
