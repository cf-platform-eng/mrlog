package main

import (
	"fmt"
	"os"

	"github.com/cf-platform-eng/mrlog"
	"github.com/cf-platform-eng/mrlog/dependency"
	"github.com/cf-platform-eng/mrlog/section"
	"github.com/jessevdk/go-flags"

	"github.com/cf-platform-eng/mrlog/version"
)

var config mrlog.Config
var parser = flags.NewParser(&config, flags.Default)

func main() {

	_, err := parser.AddCommand(
		"dependency",
		"log a dependecy",
		"log a dependency in MRL format",
		&dependency.DependencyOpt{
			Out:   os.Stdout,
			Clock: &mrlog.Clock{},
		},
	)
	if err != nil {
		fmt.Println("Could not add dependency command")
		os.Exit(1)
	}

	_, err = parser.AddCommand(
		"section-start",
		"log a section beginning",
		"log a section beginnning in MRL format",
		&section.SectionOpt{
			Section: section.Section{
				Type: "start",
			},
			Out:   os.Stdout,
			Clock: &mrlog.Clock{},
		},
	)
	if err != nil {
		fmt.Println("Could not add section command")
		os.Exit(1)
	}

	_, err = parser.AddCommand(
		"section-end",
		"log a section ending",
		"log a section endging in MRL format",
		&section.SectionOpt{
			Section: section.Section{
				Type: "end",
			},
			Out:   os.Stdout,
			Clock: &mrlog.Clock{},
		},
	)
	if err != nil {
		fmt.Println("Could not add section command")
		os.Exit(1)
	}

	_, err = parser.AddCommand(
		"version",
		"print version",
		fmt.Sprintf("print %s version", mrlog.APP_NAME),
		&version.VersionOpt{
			Out: os.Stdout,
		})
	if err != nil {
		fmt.Println("Could not add version command")
		os.Exit(1)
	}

	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
