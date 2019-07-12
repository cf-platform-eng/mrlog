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

	Scenario("logging stemcell dependencies", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a stemcell dependency")

		steps.Then("the command exits without error")
		steps.And("the result is a machine and human readable stemcell dependency log line")
	})

	Scenario("logging executable dependencies", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log an executable dependency")

		steps.Then("the command exits without error")
		steps.And("the result is a machine and human readable executable dependency log line")
	})

	Scenario("insufficient input to identify a dependency", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a dependency without sufficient data to identify it")

		steps.Then("the command exits with an error")
		steps.And("the error explains what I need to provide")
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

		define.When(`^I log a stemcell dependency`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"dependency",
				"--name",
				"light-bosh-stemcell-170.107-google-kvm-ubuntu-xenial-go_agent.tgz",
				"--hash",
				"a5387ed1ea4c61d2f7c13dfa2aa5bf6978d5e1c7",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a dependency without sufficient data to identify it$`, func() {
			logCommand := exec.Command(mrlogPath, "dependency")

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log an executable dependency$`, func() {
			logCommand := exec.Command(mrlogPath, "dependency", "--name", "marman", "--version", "2.0.1")

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

		define.Then(`^the result is a machine and human readable stemcell dependency log line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("dependency reported."))
			Eventually(commandSession.Out).Should(
				Say("Name: light-bosh-stemcell-170.107-google-kvm-ubuntu-xenial-go_agent.tgz"))
			Eventually(commandSession.Out).Should(
				Say("Hash: a5387ed1ea4c61d2f7c13dfa2aa5bf6978d5e1c7"))

			mrRE := regexp.MustCompile(`\sMRL:(.*)$`)
			machineReadableString := mrRE.FindSubmatch(commandSession.Out.Contents())

			Expect(machineReadableString).To(HaveLen(2))

			machineReadable := &struct {
				Type string `json:"type"`
				Name string `json:"name"`
				Hash string `json:"hash"`
				Time time.Time
			}{}

			err := json.Unmarshal(machineReadableString[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Name).To(Equal("light-bosh-stemcell-170.107-google-kvm-ubuntu-xenial-go_agent.tgz"))
			Expect(machineReadable.Hash).To(Equal("a5387ed1ea4c61d2f7c13dfa2aa5bf6978d5e1c7"))
			Expect(machineReadable.Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 2))
		})

		define.Then("^the result is a machine and human readable executable dependency log line$", func() {
			Eventually(commandSession.Out).Should(
				Say("dependency reported."))
			Eventually(commandSession.Out).Should(
				Say("Name: marman"))
			Eventually(commandSession.Out).Should(
				Say("Version: 2.0.1"))

			mrRE := regexp.MustCompile(`\sMRL:(.*)$`)
			machineReadableString := mrRE.FindSubmatch(commandSession.Out.Contents())

			Expect(machineReadableString).To(HaveLen(2))

			machineReadable := &struct {
				Type    string `json:"type"`
				Name    string `json:"name"`
				Version string `json:"version"`
				Time    time.Time
			}{}

			err := json.Unmarshal(machineReadableString[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Name).To(Equal("marman"))
			Expect(machineReadable.Version).To(Equal("2.0.1"))
			Expect(machineReadable.Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 2))

		})

		define.Then(`^the error explains what I need to provide$`, func() {
			Eventually(commandSession.Err).Should(
				Say("Insufficient data to identify a dependency"))
			Eventually(commandSession.Err).Should(
				Say("must use at least one of:"))
			Eventually(commandSession.Err).Should(
				Say("--name"))
			Eventually(commandSession.Err).Should(
				Say("name or filename"))
			Eventually(commandSession.Err).Should(
				Say("--version"))
			Eventually(commandSession.Err).Should(
				Say("--hash"))
		})

	})
})
