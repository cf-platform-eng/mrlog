// +build feature

package features_test

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// go:generate

var _ = Describe("log a dependency", func() {
	steps := NewSteps()

	Scenario("logging stemcell dependencies", func() {
		steps.Given("I have the mrlog binary")

		steps.When("I log a stemcell dependency")

		steps.Then("the command exits without error")
		steps.And("the result is a machine and human readable stemcell dependency log line")
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
			logCommand := exec.Command(mrlogPath, "dependency", "--filename", "light-bosh-stemcell-170.107-google-kvm-ubuntu-xenial-go_agent.tgz", "--hash", "a5387ed1ea4c61d2f7c13dfa2aa5bf6978d5e1c7")
			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.Then(`^the command exits without error$`, func() {
			Eventually(commandSession).Should(gexec.Exit(0))
		})

		define.Then(`^the result is a machine and human readable stemcell dependency log line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("dependency reported. Filename: light-bosh-stemcell-170.107-google-kvm-ubuntu-xenial-go_agent.tgz, Hash: a5387ed1ea4c61d2f7c13dfa2aa5bf6978d5e1c7"))

			fmt.Println(string(commandSession.Out.Contents()))

			mrRE := regexp.MustCompile(`\sMRL:(.*)$`)
			machineReadableString := mrRE.FindSubmatch(commandSession.Out.Contents())

			Expect(machineReadableString).To(HaveLen(2))

			machineReadable := &struct {
				Type     string `json:"type"`
				Filename string `json:"filename"`
				Hash     string `json:"hash"`
				Time     time.Time
			}{}

			err := json.Unmarshal(machineReadableString[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Filename).To(Equal("light-bosh-stemcell-170.107-google-kvm-ubuntu-xenial-go_agent.tgz"))
			Expect(machineReadable.Hash).To(Equal("a5387ed1ea4c61d2f7c13dfa2aa5bf6978d5e1c7"))
			Expect(machineReadable.Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 2))
		})
	})
})
