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

	Context("dependency with every flag", func() {
		BeforeEach(func() {
			context.Hash = "112233445566778899AABBCCDDEEFF"
			context.Version = "1.2.3"
			context.Name = "some-file.tgz"
		})

		It("logs the dependency", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("dependency reported."))
			Expect(out).To(Say("Name: some-file.tgz"))
			Expect(out).To(Say("Hash: 112233445566778899AABBCCDDEEFF"))
			Expect(out).To(Say("Version: 1.2.3"))

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(0))

			mrRE := regexp.MustCompile(`\sMRL:(.*)$`)
			machineReadableString := mrRE.FindSubmatch(output)

			Expect(machineReadableString).To(HaveLen(2))

			machineReadable := &struct {
				Type    string    `json:"type"`
				Hash    string    `json:"hash"`
				Name    string    `json:"name"`
				Version string    `json:"version"`
				Time    time.Time `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableString[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Hash).To(Equal("112233445566778899AABBCCDDEEFF"))
			Expect(machineReadable.Version).To(Equal("1.2.3"))
			Expect(machineReadable.Name).To(Equal("some-file.tgz"))
			Expect(machineReadable.Time).To(Equal(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC)))

		})
	})
	Context("dependency without every flag", func() {
		BeforeEach(func() {
			context.Name = "some-file.tgz"
		})
		It("Doesn't render missing version and hash", func() {
			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("dependency reported."))
			Expect(out).To(Say("Name: some-file.tgz"))
			Expect(out).To(Not(Say("Hash: ")))
			Expect(out).To(Not(Say("Version: ")))
			Expect(out).To(Say("MRL:"))
			Expect(out).To(Not(Say(`"Hash"`)))
			Expect(out).To(Not(Say(`"Version"`)))

			output := out.Contents()
			Expect(bytes.Count(output, []byte("\n"))).To(Equal(0))

			mrRE := regexp.MustCompile(`\sMRL:(.*)$`)
			machineReadableString := mrRE.FindSubmatch(output)

			Expect(machineReadableString).To(HaveLen(2))

			machineReadable := &struct {
				Type string    `json:"type"`
				Name string    `json:"name"`
				Time time.Time `json:"time"`
			}{}

			err := json.Unmarshal(machineReadableString[1], machineReadable)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineReadable.Type).To(Equal("dependency"))
			Expect(machineReadable.Name).To(Equal("some-file.tgz"))
			Expect(machineReadable.Time).To(Equal(time.Date(1973, 11, 29, 10, 15, 01, 00, time.UTC)))
		})
	})
})
