package dependency_test

import (
	"bytes"
	"encoding/json"
	"regexp"
	"time"

	"github.com/cf-platform-eng/mrlog/clock/clockfakes"
	"github.com/cf-platform-eng/mrlog/dependency"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Dependency", func() {
	var (
		out     *Buffer
		context *dependency.DependencyOpt
	)

	BeforeEach(func() {
		out = NewBuffer()

		clock := &clockfakes.FakeClock{}
		clock.NowReturns(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC))

		context = &dependency.DependencyOpt{
			Out:   out,
			Clock: clock,
		}
	})

	Context("dependency with invalid metadata", func() {
		BeforeEach(func() {
			context.Version = "1.2.3"
			context.Name = "some-file.tgz"
			context.Metadata = "I AM NOT JSON"
		})

		It("logs the dependency", func() {
			err := context.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid metadata"))
		})
	})

	Context("dependency with every flag", func() {
		BeforeEach(func() {
			context.Version = "1.2.3"
			context.Name = "some-file.tgz"
			context.Metadata = "{\"some-key\":\"some-value\"}"
		})

		It("logs the dependency", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("dependency: 'some-file.tgz' version '1.2.3'"))

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(1))

			mrRE := regexp.MustCompile(`\s(?m)MRL:(.*)\n`)

			machineReadableMatches := mrRE.FindSubmatch(output)

			machineReadable := &struct {
				Type     string      `json:"type"`
				Name     string      `json:"name"`
				Version  string      `json:"version"`
				Metadata interface{} `json:"metadata"`
				Time     time.Time   `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableMatches[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Version).To(Equal("1.2.3"))
			Expect(machineReadable.Name).To(Equal("some-file.tgz"))
			Expect(machineReadable.Time).To(Equal(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC)))
			Expect(machineReadable.Metadata).To(HaveKeyWithValue("some-key", "some-value"))
		})
	})
})
