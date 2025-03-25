package e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AI Explorer E2E Suite")
}

var _ = BeforeSuite(func() {
	By("Loading environment variables from .env")
	err := godotenv.Load("../.env")
	Expect(err).ToNot(HaveOccurred(), "Failed to load .env file")
})

var _ = AfterSuite(func() {
	By("Cleaning up test-generated artifacts")

	// Remove .build/output directory
	outputDir := filepath.Join("..", ".build", "output")
	err := os.RemoveAll(outputDir)
	Expect(err).ToNot(HaveOccurred(), "Failed to remove test output directory")

	// Remove ai-explorer binary
	binaryPath := filepath.Join("..", ".build", "ai-explorer")
	err = os.Remove(binaryPath)
	if err != nil && !os.IsNotExist(err) {
		// Ignore if already deleted, fail for other errors
		Expect(err).ToNot(HaveOccurred(), "Failed to remove compiled binary")
	}
})
