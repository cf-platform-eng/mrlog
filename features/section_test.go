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

	Describe("subshell sections", func() {
		Scenario("gracefully handles missing command", func() {
			steps.Given("I have the mrlog binary")

			steps.When("I log a section without a command")

			steps.Then("the command exits with 1")
			steps.And("tells me a command is required")

		})
		Scenario("reports failure to execute subcommand", func() {
			steps.Given("I have the mrlog binary")

			steps.When("I log a section with a missing subcommand")

			steps.Then("the command exits with -1")
			steps.And("the result is a machine and human readable section begin line")
			steps.And("the result contains human and machine readable result -1 section end line")
		})
		Scenario("logging a successful sub command", func() {
			steps.Given("I have the mrlog binary")

			steps.When("I log a section with a successful subcommand")

			steps.Then("the command exits without error")
			steps.And("the result is a machine and human readable section begin line")
			steps.And("the result contains output from the successful command")
			steps.And("the result contains human and machine readable successful section end line")
		})
		Scenario("logging a failed sub command", func() {
			steps.Given("I have the mrlog binary")

			steps.When("I log a section with a failed subcommand")

			steps.Then("the command exits with 2")
			steps.And("the result is a machine and human readable section begin line")
			steps.And("the result contains output from the failed command")
			steps.And("the result contains human and machine readable result 2 section end line")
		})
		Scenario("section with successful subcommand shows on-success message", func() {
			steps.Given("I have the mrlog binary")
			steps.When("I log a section with a successful subcommand and on-success/on-failure messages")
			steps.Then("the command exits without error")
			steps.And("the result is a machine and human readable section begin line")
			steps.And("the result contains output from the successful command")
			steps.And("the result contains human and machine readable successful section end line with success message")
		})
		Scenario("section with failed subcommand shows on-failed message", func() {
			steps.Given("I have the mrlog binary")
			steps.When("I log a section with a failed subcommand and on-success/on-failure messages")
			steps.Then("the command exits with 2")
			steps.And("the result is a machine and human readable section begin line")
			steps.And("the result contains output from the failed command")
			steps.And("the result contains human and machine readable successful section end line with failed message")
		})
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
				"--name",
				"test-section",
				"--result",
				"1",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section with a successful subcommand$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section",
				"--name",
				"test-section",
				"--",
				"fixtures/successful-subcommand.sh",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section with a failed subcommand$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section",
				"--name",
				"test-section",
				"--",
				"fixtures/failed-subcommand.sh",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section with a successful subcommand and on-success/on-failure messages$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section",
				"--name",
				"test-section",
				"--on-success",
				"this command was successful",
				"--on-failure",
				"this command was a failure",
				"--",
				"fixtures/successful-subcommand.sh",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section with a failed subcommand and on-success/on-failure messages$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section",
				"--name",
				"test-section",
				"--on-success",
				"this command was successful",
				"--on-failure",
				"this command was a failure",
				"--",
				"fixtures/failed-subcommand.sh",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section without a command$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section",
				"--name",
				"test-section",
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.When(`^I log a section with a missing subcommand$`, func() {
			logCommand := exec.Command(
				mrlogPath,
				"section",
				"--name",
				"test-section",
				"--",
				"fixtures/missing-subcommand.sh", // this file does not exist, keep it that way!
			)

			var err error
			commandSession, err = gexec.Start(logCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		define.Then(`^the command exits with 1$`, func() {
			Eventually(commandSession).Should(gexec.Exit(1))
		})

		define.Then(`^the command exits with 2$`, func() {
			Eventually(commandSession).Should(gexec.Exit(2))
		})

		define.Then(`^the command exits with -1$`, func() {
			Eventually(commandSession).Should(gexec.Exit(-1))
		})

		define.Then(`^the command exits without error$`, func() {
			Eventually(commandSession).Should(gexec.Exit(0))
		})

		define.Then(`^the result is a machine and human readable section begin line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-start: 'test-section'"))

			expectedMRL := mrl{
				Type: "section-start",
				Name: "test-section",
				Time: time.Now(),
			}

			matchMRL(&expectedMRL, `\s(?m)MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result is a machine and human readable section end line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: 'test-section' result: 1"))

			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "test-section",
				Result: 1,
				Time:   time.Now(),
			}

			matchMRL(&expectedMRL, `\s(?m)MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result contains human and machine readable successful section end line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: 'test-section' result: 0 "))

			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "test-section",
				Result: 0,
				Time:   time.Now(),
			}

			matchMRL(&expectedMRL, `section-end:.*MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result contains human and machine readable result -1 section end line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: 'test-section' result: -1 "))
			
			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "test-section",
				Result: -1,
				Time:   time.Now(),
			}

			matchMRL(&expectedMRL, `section-end:.*MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result contains human and machine readable result 2 section end line$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: 'test-section' result: 2 "))
	
			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "test-section",
				Result: 2,
				Time:   time.Now(),
			}

			matchMRL(&expectedMRL, `section-end:.*MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result contains human and machine readable successful section end line with success message$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: 'test-section' result: 0 message: 'this command was successful'"))
			
			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "test-section",
				Result: 0,
				Time:   time.Now(),
				Message: "this command was successful",
			}

			matchMRL(&expectedMRL, `section-end:.*MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result contains human and machine readable successful section end line with failed message$`, func() {
			Eventually(commandSession.Out).Should(
				Say("section-end: 'test-section' result: 2 message: 'this command was a failure'"))

			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "test-section",
				Result: 2,
				Time:   time.Now(),
				Message: "this command was a failure",
			}

			matchMRL(&expectedMRL, `section-end:.*MRL:(.*)\n`, commandSession.Out.Contents())
		})

		define.Then(`^the result contains output from the successful command$`, func() {
			Eventually(commandSession.Out).Should(
				Say("This is a successful command"))
		})

		define.Then(`^the result contains output from the failed command$`, func() {
			Eventually(commandSession.Out).Should(
				Say("This is a failed command"))
		})

		define.Then(`^the error telling me to provide a name$`, func() {
			Eventually(commandSession.Err).Should(
				Say("the required flag '--name' was not specified"))
		})

		define.Then(`tells me a command is required$`, func() {
			Eventually(commandSession.Err).Should(
				Say("the section subcommand requires a command parameter '-- <command> ...'"))
		})
	})
})

type mrl struct {
	Type    string    `json:"type"`
	Name    string    `json:"name"`
	Result  int       `json:"result"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func matchMRL(expectedMRL *mrl, regex string, output []byte) {
	mrRE := regexp.MustCompile(regex)

	machineReadableMatches := mrRE.FindSubmatch(output)

	Expect(len(machineReadableMatches)).To(Equal(2))
	var machineReadable mrl

	err := json.Unmarshal(machineReadableMatches[1], &machineReadable)
	Expect(err).NotTo(HaveOccurred())

	Expect(machineReadable.Type).To(Equal(expectedMRL.Type))
	Expect(machineReadable.Name).To(Equal(expectedMRL.Name))
	Expect(machineReadable.Result).To(Equal(expectedMRL.Result))
	Expect(machineReadable.Time.Unix()).To(BeNumerically("~", expectedMRL.Time.Unix(), 2))
	Expect(machineReadable.Message).To(Equal(expectedMRL.Message))
}
