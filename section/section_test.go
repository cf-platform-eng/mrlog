package section_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	os_exec "os/exec"

	"github.com/cf-platform-eng/mrlog/section"
	"github.com/cf-platform-eng/mrlog/section/sectionfakes"
	"github.com/cf-platform-eng/mrlog/exec/execfakes"
	"github.com/cf-platform-eng/mrlog/clock/clockfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

//go:generate counterfeiter io.Writer

var _ = Describe("Section", func() {
	var (
		out     *Buffer
		context *section.SectionOpt
		cmd *execfakes.FakeCmd
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
			Exec: exec,
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

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(1))

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)

			machineReadableMatches := mrRE.FindSubmatch(output)

			machineReadable := &struct {
				Type string    `json:"type"`
				Name string    `json:"name"`
				Time time.Time `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("section-start"))
			Expect(machineReadable.Name).To(Equal("install"))
			Expect(machineReadable.Time).To(Equal(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC)))
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

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(1))

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)

			machineReadableMatches := mrRE.FindSubmatch(output)

			machineReadable := &struct {
				Type   string    `json:"type"`
				Name   string    `json:"name"`
				Result int       `json:"result"`
				Time   time.Time `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("section-end"))
			Expect(machineReadable.Name).To(Equal("install"))
			Expect(machineReadable.Result).To(Equal(1))
			Expect(machineReadable.Time).To(Equal(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC)))

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
			err := context.Execute([]string{"command"})
			Expect(err).ToNot(HaveOccurred())
		})

		Context("reports failure to execute subcommand", func() {
			BeforeEach(func() {
				cmd.RunReturns(fmt.Errorf("command failed"))
			})

			It("fails when command run fails", func() {
				err := context.Execute([]string{"command"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("command failed"))
			})
		})
	})
})
