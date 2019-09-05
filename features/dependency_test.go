// +build feature

package features_test

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"time"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("log a dependency", func() {
	steps := NewSteps()

	Scenario("logging a dependency", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a dependency")

		steps.Then("the command exits without error")
		steps.And("the result contains a human readable log")
		steps.And("the result contains a machine readable log")
	})

	Scenario("logging a dependency with metadata", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a dependency with metadata")

		steps.Then("the command exits without error")
		steps.And("the result contains a human readable log")
		steps.And("the result contains a machine readable log")
		steps.And("the machine readable dependency log contains provided metadata")
	})

	Scenario("logging a dependency without a name", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a dependency without a name")

		steps.Then("the command exits with an error")
	})

	Scenario("logging a dependency without a version", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a dependency without a version")

		steps.Then("the command exits with an error")
	})

	steps.Define(func(define Definitions) {
		var (
			commandSession *gexec.Session
			mrlogPath      string
		)

		define.Given(`^I have the mrlog binary$`, func() {
			var err error
			mrlogPath, err = gexec.Build("github.com/cf-platform-eng/mrlog/cmd/mrlog")
			Expect(err).NotTo(HaveOccurred())
		}, func() {
			gexec.CleanupBuildArtifacts()
		})

		define.When(`^I log a dependency$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"dependency",
				"--name",
				"marman",
				"--version",
				"1.2.3",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a dependency without a name$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"dependency",
				"--version",
				"1.2.3",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a dependency without a version$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"dependency",
				"--name",
				"marman",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a dependency with metadata$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"dependency",
				"--name",
				"marman",
				"--version",
				"1.2.3",
				"--metadata",
				"{\"some-key\":\"some-value\"}",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.Then(`^the command exits without error$`, func() {
			Eventually(commandSession).Should(gexec.Exit(0))
		})

		define.Then(`^the command exits with an error$`, func() {
			Eventually(commandSession).Should(gexec.Exit(1))
		})

		define.Then(`^the result contains a human readable log$`, func() {
			Eventually(commandSession.Out).Should(
				Say("dependency: 'marman' version '1.2.3'"))
		})

		define.Then(`^the result contains a machine readable log$`, func() {
			contents := commandSession.Out.Contents()

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)
			Expect(mrRE.Match(contents)).To(BeTrue())

			machineReadableMatches := mrRE.FindSubmatch(contents)

			machineReadable := &struct {
				Type     string      `json:"type"`
				Name     string      `json:"name"`
				Version  string      `json:"version"`
				Metadata interface{} `json:"metadata"`
				Time     time.Time
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Name).To(Equal("marman"))
			Expect(machineReadable.Version).To(Equal("1.2.3"))
			Expect(machineReadable.Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 2))
		})

		define.Then(`^the machine readable dependency log contains provided metadata$`, func() {
			contents := commandSession.Out.Contents()

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)
			Expect(mrRE.Match(contents)).To(BeTrue())

			machineReadableMatches := mrRE.FindSubmatch(contents)

			machineReadable := &struct {
				Type     string      `json:"type"`
				Name     string      `json:"name"`
				Version  string      `json:"version"`
				Metadata interface{} `json:"metadata"`
				Time     time.Time
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Metadata).To(HaveKeyWithValue("some-key", "some-value"))
		})

		define.Then(`^the error telling me to provide a name$`, func() {
			Eventually(commandSession.Err).Should(
				Say("the required flag `--name' was not specified"))
		})
	})
})
