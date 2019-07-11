package main

import (
	"fmt"
	"os"

	"github.com/cf-platform-eng/mrlog"

	"github.com/jessevdk/go-flags"

	"github.com/cf-platform-eng/mrlog/version"
)

var config mrlog.Config
var parser = flags.NewParser(&config, flags.Default)

func main() {
	_, err := parser.AddCommand(
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
