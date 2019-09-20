package dependency

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
	"github.com/pkg/errors"
)

type Identities struct {
	Name           string `long:"name" description:"name of the dependency, if it has one" required:"true"`
	Version        string `long:"version" description:"version string for the dependency, if it has one" required:"true"`
	Metadata       string `long:"metadata" description:"optionally provide metadata for this dependency"`
	DependencyType string `long:"type" description:"type of dependency"`
}

type DependencyOpt struct {
	Identities
	Out   io.Writer
	Clock clock.Clock
}

func (opts *DependencyOpt) Execute(args []string) error {
	dependency := "dependency"
	if opts.DependencyType != "" {
		dependency = fmt.Sprintf("%s dependency", opts.DependencyType)
	}
	humanReadable := fmt.Sprintf("%s: "+
		"'%s' version '%s'",
		dependency,
		opts.Name,
		opts.Version)

	_, err := fmt.Fprint(opts.Out, humanReadable)
	if err != nil { // !branch-not-tested
		return err
	}

	machineLog := &mrl.MachineReadableLog{
		Type:     dependency,
		Version:  opts.Version,
		Name:     opts.Name,
		Metadata: opts.Metadata,
		Time:     opts.Clock.Now(),
	}

	if opts.Metadata != "" {
		err = json.Unmarshal([]byte(opts.Metadata), &machineLog.Metadata)
		if err != nil {
			return errors.Wrap(err, "invalid metadata")
		}
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
