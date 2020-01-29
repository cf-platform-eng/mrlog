package section_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/cf-platform-eng/mrlog/clock/clockfakes"
	"github.com/cf-platform-eng/mrlog/exec/execfakes"
	"github.com/cf-platform-eng/mrlog/section"
	"github.com/cf-platform-eng/mrlog/section/sectionfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

//go:generate counterfeiter io.Writer

var _ = Describe("Section", func() {
	var (
		out     *Buffer
		context *section.SectionOpt
		cmd     *execfakes.FakeCmd
	)

	BeforeEach(func() {
		out = NewBuffer()

		clock := &clockfakes.FakeClock{}
		clock.NowReturns(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC))
		exec := &execfakes.FakeExec{}
		cmd = &execfakes.FakeCmd{}
		exec.CommandReturns(cmd)

		context = &section.SectionOpt{
			Out:   out,
			Clock: clock,
			Exec:  exec,
		}
	})

	Context("section with a name", func() {
		BeforeEach(func() {
			context.Name = "install"
			context.Type = "start"
		})

		It("logs the section", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("section-start: 'install'"))

			expectedMRL := mrl{
				Type: "section-start",
				Name: "install",
				Time: time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
			}
			matchMRL(&expectedMRL, `\s(?m)MRL:(.*)\n`, out.Contents())
		})
	})

	Context("section without a name", func() {
		BeforeEach(func() {
			context.Name = ""
			context.Type = "start"
		})

		It("return an error", func() {
			err := context.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("missing section name"))
		})
	})

	Context("invalid section type", func() {
		BeforeEach(func() {
			context.Name = "pete"
			context.Type = "coding"
		})

		It("return an error", func() {
			err := context.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid section type argument"))
		})
	})

	Context("section end", func() {
		BeforeEach(func() {
			context.Result = 1
			context.Type = "end"
			context.Name = "install"
		})

		It("logs the section end", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("section-end: 'install' result: 1"))

			expectedMRL := mrl{
				Type:   "section-end",
				Name:   "install",
				Time:   time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
				Result: 1,
			}

			matchMRL(&expectedMRL, `\s(?m)MRL:(.*)\n`, out.Contents())
		})
	})

	Context("section output fails", func() {
		var output *sectionfakes.FakeWriter

		BeforeEach(func() {
			output = &sectionfakes.FakeWriter{}
			output.WriteReturns(0, errors.New("write-error"))
			context.Out = output
			context.Name = "install"
			context.Type = "start"
		})

		It("logs the section", func() {
			err := context.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("write-error"))
			Expect(err.Error()).To(ContainSubstring("failed to write"))
		})
	})

	Context("section subcommand", func() {
		BeforeEach(func() {
			context.Type = "section"
			context.Name = "install"
		})

		It("fails given no command", func() {
			err := context.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("the section subcommand requires a command parameter '-- <command> ...'"))
		})

		It("succeeds given command", func() {
			Expect(context.Execute([]string{"command"})).To(Succeed())
			Expect(out).To(Say("section-start: 'install'"))
			Expect(out).To(Say("section-end: 'install' result: 0"))

			startMRL := mrl{
				Type: "section-start",
				Name: "install",
				Time: time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
			}
			matchMRL(&startMRL, `section-start:.*MRL:(.*)`, out.Contents())
			endMRL := mrl{
				Type:   "section-end",
				Name:   "install",
				Time:   time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
				Result: 0,
			}
			matchMRL(&endMRL, `section-end:.*MRL:(.*)`, out.Contents())
		})

		Context("reports failure to execute subcommand", func() {
			BeforeEach(func() {
				cmd.RunReturns(fmt.Errorf("command failed"))
			})

			It("fails when command run fails", func() {
				err := context.Execute([]string{"command"})
				Expect(err).To(HaveOccurred())
				Expect(out).To(Say("section-start: 'install'"))
				Expect(out).To(Say("command failed"))
				Expect(out).To(Say("section-end: 'install' result: -1"))

				startMRL := mrl{
					Type: "section-start",
					Name: "install",
					Time: time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
				}
				matchMRL(&startMRL, `section-start:.*MRL:(.*)`, out.Contents())
				endMRL := mrl{
					Type:   "section-end",
					Name:   "install",
					Time:   time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
					Result: -1,
				}
				matchMRL(&endMRL, `section-end:.*MRL:(.*)`, out.Contents())
			})
		})

		Context("success/failure messages", func() {
			BeforeEach(func() {
				context.Type = "section"
				context.Name = "messages"
				context.OnSuccess = "successful"
				context.OnFailure = "failure"
			})
			Context("success", func() {
				It("prints success message", func() {
					Expect(context.Execute([]string{"command"})).To(Succeed())
					Expect(out).To(Say("section-end: 'messages' result: 0 message: 'successful'"))

					expectedMRL := mrl{
						Type:    "section-end",
						Name:    "messages",
						Result:  0,
						Time:    time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
						Message: "successful",
					}
					matchMRL(&expectedMRL, `section-end:.*MRL:(.*)`, out.Contents())
				})
			})
			Context("failure", func() {
				BeforeEach(func() {
					cmd.RunReturns(fmt.Errorf("command failed"))
				})
				It("prints failure message", func() {
					Expect(context.Execute([]string{"command"})).NotTo(Succeed())
					Expect(out).To(Say("section-end: 'messages' result: -1 message: 'failure'"))

					expectedMRL := mrl{
						Type:    "section-end",
						Name:    "messages",
						Result:  -1,
						Time:    time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC),
						Message: "failure",
					}
					matchMRL(&expectedMRL, `section-end:.*MRL:(.*)`, out.Contents())
				})
			})
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
	Expect(machineReadable.Time).To(Equal(expectedMRL.Time))
	Expect(machineReadable.Message).To(Equal(expectedMRL.Message))
}
