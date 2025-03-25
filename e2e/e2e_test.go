package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// --- Configurable Constants ---
const (
	binaryName      = "ai-explorer"
	topic           = "git"
	openai_model    = "gpt-4o"
	ollama_model    = "phi4"
	openai_provider = "openai"
	ollama_provider = "ollama"
	temperature     = "0.7"
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
func newTestPaths(topic string, suffix string) *TestPaths {
	root := ".."
	output := filepath.Join(root, ".build", "output", "test_"+topic+"_"+suffix)

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

func TestAIExplorerCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AI Explorer CLI E2E Suite")
}

var _ = Describe("AI Explorer CLI (E2E)", func() {
	// Parallelize the entire describe block
	GinkgoParallelProcess()

	// Prompt Generation Test
	Describe("Prompt Generation", func() {
		It("Should generate a prompt file", func() {
			paths := newTestPaths(topic, "prompt")
			err := os.MkdirAll(paths.OutputDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "Failed to create output directory")

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

	// LLM Command Tests
	Describe("LLM Commands", func() {
		// OpenAI LLM Test
		It("Should return a valid OpenAI model response", func() {
			paths := newTestPaths(topic, "openai_llm")
			err := os.MkdirAll(paths.OutputDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "Failed to create output directory")

			// First generate a prompt
			promptOutput := filepath.Join(paths.OutputDir, "prompt.txt")
			output, err := runCommand(paths,
				"prompt",
				"--topic", topic,
				"--template", paths.TemplatePath,
				"--config", paths.ConfigPath,
				"--output", promptOutput,
			)
			Expect(err).ToNot(HaveOccurred(), "Prompt generation failed:\n%s", string(output))

			// Then run LLM command
			output, err = runCommand(paths,
				"llm",
				"--provider", openai_provider,
				"--model", openai_model,
				"--prompt", promptOutput,
				"--temperature", temperature,
			)

			Expect(err).ToNot(HaveOccurred(), "LLM command failed:\n%s", string(output))
			Expect(string(output)).To(ContainSubstring(topic))
		})

		// Ollama LLM Test
		It("Should return a valid Ollama model response", func() {
			paths := newTestPaths(topic, "ollama_llm")
			err := os.MkdirAll(paths.OutputDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "Failed to create output directory")

			// First generate a prompt
			promptOutput := filepath.Join(paths.OutputDir, "prompt.txt")
			output, err := runCommand(paths,
				"prompt",
				"--topic", topic,
				"--template", paths.TemplatePath,
				"--config", paths.ConfigPath,
				"--output", promptOutput,
			)
			Expect(err).ToNot(HaveOccurred(), "Prompt generation failed:\n%s", string(output))

			// Then run LLM command
			output, err = runCommand(paths,
				"llm",
				"--provider", ollama_provider,
				"--model", ollama_model,
				"--prompt", promptOutput,
				"--temperature", temperature,
			)

			Expect(err).ToNot(HaveOccurred(), "LLM command failed:\n%s", string(output))
			Expect(string(output)).To(ContainSubstring(topic))
		})
	})

	// Chat Command Tests
	Describe("Chat Commands", func() {
		// OpenAI Chat Test
		It("Should generate a prompt and get an OpenAI response", func() {
			paths := newTestPaths(topic, "openai_chat")
			err := os.MkdirAll(paths.OutputDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "Failed to create output directory")

			output, err := runCommand(paths,
				"chat",
				"--topic", topic,
				"--provider", openai_provider,
				"--model", openai_model,
			)

			Expect(err).ToNot(HaveOccurred(), "Chat command failed:\n%s", string(output))
			Expect(string(output)).To(ContainSubstring(topic))
		})

		// Ollama Chat Test
		It("Should generate a prompt and get an Ollama response", func() {
			paths := newTestPaths(topic, "ollama_chat")
			err := os.MkdirAll(paths.OutputDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "Failed to create output directory")

			output, err := runCommand(paths,
				"chat",
				"--topic", topic,
				"--provider", ollama_provider,
				"--model", ollama_model,
			)

			Expect(err).ToNot(HaveOccurred(), "Chat command failed:\n%s", string(output))
			Expect(string(output)).To(ContainSubstring(topic))
		})
	})
})
