package version

import (
	"fmt"
	"io"

	"github.com/cf-platform-eng/mrlog"
)


type VersionOpt struct {
	Out io.Writer
}

var Version = "dev"

func (opts *VersionOpt) Execute(args []string) error {
	_, err := fmt.Fprintf(opts.Out, "%s version: %s\n", mrlog.APP_NAME, Version)
	return err
}

