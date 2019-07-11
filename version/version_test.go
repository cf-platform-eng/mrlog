package version_test

import (
	"github.com/cf-platform-eng/mrlog/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Version", func() {
	Context("version is set", func() {
		var (
			originalVersion string
		)

		BeforeEach(func() {
			originalVersion = version.Version
			version.Version = "9.9.9"
		})
		AfterEach(func() {
			version.Version = originalVersion
		})

		It("prints the version", func() {
			out := NewBuffer()
			context := version.VersionOpt{
				Out: out,
			}

			Expect(context.Execute([]string{})).To(Succeed())
			Expect(out).To(Say("mrlog version: 9.9.9"))
		})

	})
})
