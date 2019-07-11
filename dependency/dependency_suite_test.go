package dependency_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVersion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependency Suite")
}
