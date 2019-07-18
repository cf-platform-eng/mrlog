package section

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type Section struct {
	Type   string
	Name   string `long:"name" description:"name of the section"`
	Result int    `long:"result" description:"result code for section"`
}

type SectionOpt struct {
	Section
	Out   io.Writer
	Clock clock.Clock
}

func (opts *SectionOpt) Execute(args []string) error {
	var humanReadable string
	if opts.Type == "start" {
		humanReadable = fmt.Sprintf("section-%s: '%s'",
			opts.Type,
			opts.Name)

	} else if opts.Type == "end" {
		humanReadable = fmt.Sprintf("section-%s: '%s' result: %d",
			opts.Type,
			opts.Name,
			opts.Result)
	} else {
		return errors.New("invalid section type argument")
	}

	if opts.Name == "" {
		return errors.New("missing section name")
	}

	_, err := fmt.Fprint(opts.Out, humanReadable)
	if err != nil {
		return errors.Wrap(err, "failed to write")
	}

	machineLog := &mrl.MachineReadableLog{
		Name:   opts.Name,
		Type:   fmt.Sprintf("section-%s", opts.Type),
		Result: opts.Result,
		Time:   opts.Clock.Now(),
	}

	machineLogJSON, err := json.Marshal(machineLog)
	if err != nil { // !branch-not-tested
		return err
	}

	_, err = fmt.Fprintf(opts.Out, " MRL:%s\n", string(machineLogJSON))
	if err != nil {
		return errors.Wrap(err, "failed to write")
	}

	return nil

}
