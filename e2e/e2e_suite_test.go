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

func createDir(path string) {
	Expect(os.MkdirAll(path, os.ModePerm)).
		To(Succeed(), "Failed to create directory: %s", path)
	GinkgoWriter.Printf("Created directory: %s\n", path)
}

func cleanupPath(path string, mustExist bool) {
	err := os.RemoveAll(path)
	if mustExist || (!mustExist && !os.IsNotExist(err)) {
		Expect(err).ToNot(HaveOccurred(), "Failed to remove: %s", path)
	}
}

var _ = BeforeSuite(func() {
	By("Loading environment variables from .env")
	Expect(godotenv.Load("../.env")).To(Succeed(), "Failed to load .env file")

	By("Creating output directory")
	createDir(filepath.Join("..", ".build", "output"))
})

var _ = AfterSuite(func() {
	By("Cleaning up test-generated artifacts")
	cleanupPath(filepath.Join("..", ".build", "output"), false)
	cleanupPath(filepath.Join("..", ".build", "ai-explorer"), false)
})
