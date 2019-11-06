package section

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type Section struct {
	Type   string
	Name   string `long:"name" description:"name of the section"`
	Result int    `long:"result" description:"exitCode code for section"`
}

type SectionOpt struct {
	Section
	Out   io.Writer
	Clock clock.Clock
}

type SectionError struct {
	Retval int
	Err error
}

func writeSection(sectionType, sectionName string, exitCode int, time time.Time, out io.Writer) error {
	var humanReadable string
	if sectionType == "start" {
		humanReadable = fmt.Sprintf("section-%s: '%s'",
			sectionType,
			sectionName)
	} else if sectionType == "end" {
		humanReadable = fmt.Sprintf("section-%s: '%s' result: %d",
			sectionType,
			sectionName,
			exitCode)
	} else {
		return errors.New("invalid section type argument")
	}

	_, err := fmt.Fprint(out, humanReadable)
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	machineLog := &mrl.MachineReadableLog{
		Name:   sectionName,
		Type:   fmt.Sprintf("section-%s", sectionType),
		Result: exitCode,
		Time:   time,
	}

	machineLogJSON, err := json.Marshal(machineLog)
	if err != nil { // !branch-not-tested
		return err
	}

	_, err = fmt.Fprintf(out, " MRL:%s\n", string(machineLogJSON))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func (e *SectionError) Unwrap() error { return e.Err }
func (e *SectionError) Error() string { return fmt.Sprintf("Section subcommand failed with %d: %s", e.Retval, e.Err)}

func (opts *SectionOpt) Execute(args []string) error {
	if opts.Name == "" {
		return errors.New("missing section name")
	}

	if opts.Type == "section" {
		if len(args) == 0 {
			return errors.New("the section subcommand requires a command parameter '-- <command> ...'")
		}

		if err := writeSection("start", opts.Name, 0, opts.Clock.Now(), opts.Out); err != nil {
			return err
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = opts.Out
		cmd.Stderr = opts.Out
		err := cmd.Run()

		exitCode := 0

		var sectionError *SectionError

		if err != nil {
			var e *exec.ExitError
			if errors.As(err, &e) {
				exitCode = e.ExitCode()
			} else {
				exitCode = -1
			}
			sectionError = &SectionError{exitCode, err}
		}
		
		err = writeSection("end", opts.Name, exitCode, opts.Clock.Now(), opts.Out)

		if sectionError != nil {
			// returning sectionError to propagate subcommand result code
			return sectionError
		}
		return err
	}

	return writeSection(opts.Type, opts.Name, opts.Result, opts.Clock.Now(), opts.Out)
}
