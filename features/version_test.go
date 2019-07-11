// +build feature

package features_test

import (
	"os/exec"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Report version", func() {
	steps := NewSteps()

	Scenario("version command reports version", func() {
		steps.Given("the mrlog command is built with a version")

		steps.When("version subcommand is run")

		steps.Then("the command exits without error")
		steps.And("the result is the version")
	})

	steps.Define(func(define Definitions) {
		var (
			mrlogPath      string
			commandSession *gexec.Session
		)

		define.Given(`^the mrlog command is built with a version$`, func() {
			var err error
			mrlogPath, err = gexec.Build(
				"github.com/cf-platform-eng/mrlog/cmd/mrlog",
				"-ldflags",
				"-X github.com/cf-platform-eng/mrlog/version.Version=1.0.1",
			)
			Expect(err).NotTo(HaveOccurred())
		}, func() {
			gexec.CleanupBuildArtifacts()
		})

		define.When(`^version subcommand is run$`, func() {
			versionCommand := exec.Command(mrlogPath, "version")
			var err error
			commandSession, err = gexec.Start(versionCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.Then(`^the command exits without error$`, func() {
			Eventually(commandSession).Should(gexec.Exit(0))
		})

		define.Then(`the result is the version`, func() {
			Eventually(commandSession.Out).Should(Say("mrlog version: 1.0.1"))
		})
	})
})
