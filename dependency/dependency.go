package dependency

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type Identities struct {
	Hash    string `long:"hash" description:"hash sum of the dependency, if it has one"`
	Name    string `long:"name" description:"name of the dependency, if it has one" required:"true"`
	Version string `long:"version" description:"version string for the dependency, if it has one"`
}

type DependencyOpt struct {
	Identities
	Out   io.Writer
	Clock clock.Clock
}

func (opts *DependencyOpt) humanReadableOutput() string {
	humanLog := fmt.Sprintf("dependency reported. "+
		"Name: %s",
		opts.Name)

	if opts.Hash != "" {
		humanLog += fmt.Sprintf(" Hash: %s ", opts.Hash)
	}
	if opts.Version != "" {
		humanLog += fmt.Sprintf(" Version: %s", opts.Version)
	}

	return humanLog
}

func (opts *DependencyOpt) Execute(args []string) error {

	_, err := fmt.Fprint(opts.Out, opts.humanReadableOutput())
	if err != nil { // !branch-not-tested
		return err
	}

	machineLog := &mrl.MachineReadableLog{
		Type:    "dependency",
		Hash:    opts.Hash,
		Version: opts.Version,
		Name:    opts.Name,
		Time:    opts.Clock.Now(),
	}

	machineLogJSON, err := json.Marshal(machineLog)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(opts.Out, " MRL:"+string(machineLogJSON))
	if err != nil { // !branch-not-tested
		return err
	}

	return nil
}
