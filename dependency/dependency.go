package dependency

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type Identities struct {
	Name    string `long:"name" description:"name of the dependency, if it has one" required:"true"`
	Version string `long:"version" description:"version string for the dependency, if it has one" required:"true"`
}

type DependencyOpt struct {
	Identities
	Out   io.Writer
	Clock clock.Clock
}

func (opts *DependencyOpt) Execute(args []string) error {

	humanReadable := fmt.Sprintf("dependency: "+
		"'%s' version '%s'",
		opts.Name,
		opts.Version)

	_, err := fmt.Fprint(opts.Out, humanReadable)
	if err != nil { // !branch-not-tested
		return err
	}

	machineLog := &mrl.MachineReadableLog{
		Type:    "dependency",
		Version: opts.Version,
		Name:    opts.Name,
		Time:    opts.Clock.Now(),
	}

	machineLogJSON, err := json.Marshal(machineLog)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(opts.Out, " MRL:%s\n", string(machineLogJSON))
	if err != nil { // !branch-not-tested
		return err
	}

	return nil
}
