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

var _ = Describe("log section boundaries", func() {
	steps := NewSteps()

	Scenario("logging beginning of a section", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a section start")

		steps.Then("the command exits without error")
		steps.And("the result is a machine and human readable section begin line")
	})

	Scenario("logging end of a section", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a section end")

		steps.Then("the command exits without error")
		steps.And("the result is a machine and human readable section end line")
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

		define.When(`^I log a section start`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section-start",
				"--name",
				"test-section",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section end`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section-end",
				"--result",
				"1",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.Then(`^the command exits without error$`, func() {
			Eventually(commandSession).Should(gexec.Exit(0))
		})

		define.Then(`^the result is a machine and human readable section begin line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-start: 'test-section'"))

			contents := commandSession.Out.Contents()

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)
			Expect(mrRE.Match(contents)).To(BeTrue())

			machineReadableMatches := mrRE.FindSubmatch(contents)

			machineReadable := &struct {
				Type string `json:"type"`
				Name string `json:"name"`
				Time time.Time
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("section-start"))
			Expect(machineReadable.Name).To(Equal("test-section"))
			Expect(machineReadable.Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 2))
		})

		define.Then(`^the result is a machine and human readable section end line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: result: 1"))

			contents := commandSession.Out.Contents()

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)
			Expect(mrRE.Match(contents)).To(BeTrue())

			machineReadableMatches := mrRE.FindSubmatch(contents)

			machineReadable := &struct {
				Type   string `json:"type"`
				Result int    `json:"result"`
				Time   time.Time
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("section-end"))
			Expect(machineReadable.Result).To(Equal(1))
			Expect(machineReadable.Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 2))
		})

		define.Then(`^the error telling me to provide a name$`, func() {
			Eventually(commandSession.Err).Should(
				Say("the required flag `--name' was not specified"))
		})

	})
})
