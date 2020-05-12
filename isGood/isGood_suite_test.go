package isGood_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIsGood(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IsGood Suite")
}
