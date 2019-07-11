package dependency

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type DependencyOpt struct {
	Filename string `long:"filename" description:"name of the dependency if it is a file"`
	Hash     string `long:"hash" description:"hash sum of the dependency if it has one"`
	Out      io.Writer
	Clock    clock.Clock
}

func (opts *DependencyOpt) Execute(args []string) error {
	humanLog := fmt.Sprintf("dependency reported. Filename: %s, Hash: %s", opts.Filename, opts.Hash)

	_, err := fmt.Fprint(opts.Out, humanLog)
	if err != nil { // !branch-not-tested
		return err
	}

	machineLog := &mrl.MachineReadableLog{
		Type:     "dependency",
		Filename: opts.Filename,
		Hash:     opts.Hash,
		Time:     opts.Clock.Now(),
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
