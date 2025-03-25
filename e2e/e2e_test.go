package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// --- Constants ---
const (
	binaryName     = "ai-explorer"
	topic          = "git"
	openaiModel    = "gpt-4o"
	ollamaModel    = "phi4"
	openaiProvider = "openai"
	ollamaProvider = "ollama"
	temperature    = "0.7"
	rootDir        = ".."
)

// --- Paths ---
type TestPaths struct {
	RootDir      string
	BinPath      string
	TemplatePath string
	ConfigPath   string
	OutputDir    string
	PromptOutput string
}

func newTestPaths(topic, suffix string) *TestPaths {
	output := filepath.Join(".build", "output", "test_"+topic+"_"+suffix)

	return &TestPaths{
		RootDir:      rootDir,
		BinPath:      filepath.Join(".build", binaryName),
		TemplatePath: filepath.Join("resources", "templates/topic.yaml"),
		ConfigPath:   filepath.Join("resources", "configs", topic+".yaml"),
		OutputDir:    output,
		PromptOutput: filepath.Join(output, "prompt.txt"),
	}
}

// --- Utilities ---
func runCommand(paths *TestPaths, args ...string) ([]byte, error) {
	GinkgoWriter.Println("Running command:", paths.BinPath, args)
	cmd := exec.Command(paths.BinPath, args...)
	cmd.Dir = paths.RootDir
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	return output, err
}

func generatePrompt(paths *TestPaths) {
	output, err := runCommand(paths,
		"prompt",
		"--topic", topic,
		"--template", paths.TemplatePath,
		"--config", paths.ConfigPath,
		"--output", paths.PromptOutput,
	)
	Expect(err).ToNot(HaveOccurred(), "Prompt generation command failed")
	GinkgoWriter.Printf("Prompt generation output:\n%s\n", string(output))

	fullPromptPath := filepath.Join(paths.RootDir, paths.PromptOutput)
	Expect(fullPromptPath).To(BeAnExistingFile(),
		"Expected prompt output file to exist at: %s", fullPromptPath)
}

func runLLMCommand(paths *TestPaths, provider, model string) {
	output, err := runCommand(paths,
		"llm",
		"--provider", provider,
		"--model", model,
		"--prompt", paths.PromptOutput,
		"--temperature", temperature,
	)
	Expect(err).ToNot(HaveOccurred(), "LLM command failed:\n%s", string(output))
	Expect(string(output)).To(ContainSubstring(topic))
}

func runChatCommand(paths *TestPaths, provider, model string) {
	output, err := runCommand(paths,
		"chat",
		"--topic", topic,
		"--provider", provider,
		"--model", model,
	)
	Expect(err).ToNot(HaveOccurred(), "Chat command failed:\n%s", string(output))
	Expect(string(output)).To(ContainSubstring(topic))
}

// --- Test Runner ---
func TestAIExplorerCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AI Explorer CLI E2E Suite")
}

// --- Test Specs ---
var _ = Describe("AI Explorer CLI (E2E)", func() {
	GinkgoParallelProcess()
	defer GinkgoRecover()

	Describe("Prompt Generation", func() {
		It("Should generate a prompt file", func() {
			paths := newTestPaths(topic, "prompt")
			generatePrompt(paths)
		})
	})

	Describe("LLM Commands", func() {
		It("Should return a valid OpenAI model response", func() {
			paths := newTestPaths(topic, "openai_llm")
			generatePrompt(paths)
			runLLMCommand(paths, openaiProvider, openaiModel)
		})

		It("Should return a valid Ollama model response", func() {
			paths := newTestPaths(topic, "ollama_llm")
			generatePrompt(paths)
			runLLMCommand(paths, ollamaProvider, ollamaModel)
		})
	})

	Describe("Chat Commands", func() {
		It("Should generate a prompt and get an OpenAI response", func() {
			paths := newTestPaths(topic, "openai_chat")
			runChatCommand(paths, openaiProvider, openaiModel)
		})

		It("Should generate a prompt and get an Ollama response", func() {
			paths := newTestPaths(topic, "ollama_chat")
			runChatCommand(paths, ollamaProvider, ollamaModel)
		})
	})
})
