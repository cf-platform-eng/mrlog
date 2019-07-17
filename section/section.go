package section

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type Section struct {
	Type string
	Name string `long:"name" description:"name of the section" required:"true"`
}

type SectionOpt struct {
	Section
	Out io.Writer
	Clock clock.Clock
}

func (opts *SectionOpt) Execute(args []string) error {
	humanReadable := fmt.Sprintf("section-%s: '%s'",
		opts.Type,
		opts.Name)

	_, err := fmt.Fprint(opts.Out, humanReadable)
	if err != nil {
		return errors.Wrap(err, "failed to write")
	}

	machineLog := &mrl.MachineReadableLog{
		Type: fmt.Sprintf("section-%s", opts.Type),
		Name: opts.Name,
		Time: opts.Clock.Now(),
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