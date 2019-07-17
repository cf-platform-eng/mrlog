package section_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/cf-platform-eng/mrlog/section"
	"github.com/cf-platform-eng/mrlog/section/sectionfakes"

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
	)

	BeforeEach(func() {
		out = NewBuffer()

		clock := &clockfakes.FakeClock{}
		clock.NowReturns(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC))

		context = &section.SectionOpt{
			Out:   out,
			Clock: clock,
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

	Context("section end", func() {
		BeforeEach(func() {
			context.Result = 1
			context.Type = "end"
		})

		It("logs the section end", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("section-end: result: 1"))

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(1))

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)

			machineReadableMatches := mrRE.FindSubmatch(output)

			machineReadable := &struct {
				Type   string    `json:"type"`
				Result int       `json:"result"`
				Time   time.Time `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("section-end"))
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
})
