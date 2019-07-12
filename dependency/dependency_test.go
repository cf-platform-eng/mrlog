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

	Context("dependency without identifying information", func() {
		It("logs the dependency", func() {
			err := context.Execute([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(dependency.InsufficientMessage))
		})
	})

	Context("dependency with a filename and hash", func() {
		BeforeEach(func() {
			context.Filename = "my-file"
			context.Hash = "112233445566778899AABBCCDDEEFF"
		})

		It("logs the dependency", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("dependency reported. Filename: my-file, Hash: 112233445566778899AABBCCDDEEFF"))

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(0))

			mrRE := regexp.MustCompile(`\sMRL:(.*)$`)
			machineReadableString := mrRE.FindSubmatch(output)

			Expect(machineReadableString).To(HaveLen(2))

			machineReadable := &struct {
				Type     string    `json:"type"`
				Filename string    `json:"filename"`
				Hash     string    `json:"hash"`
				Time     time.Time `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableString[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Filename).To(Equal("my-file"))
			Expect(machineReadable.Hash).To(Equal("112233445566778899AABBCCDDEEFF"))
			Expect(machineReadable.Time).To(Equal(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC)))

		})
	})
})
