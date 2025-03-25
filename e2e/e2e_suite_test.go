package e2e_test

import (
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
