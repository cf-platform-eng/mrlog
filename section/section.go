package section

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	os_exec "os/exec"

	"github.com/cf-platform-eng/mrlog/clock"
	"github.com/cf-platform-eng/mrlog/exec"
	"github.com/cf-platform-eng/mrlog/mrl"
)

type Section struct {
	Type      string
	Name      string `long:"name" description:"name of the section"`
	Result    int    `long:"result" description:"exitCode code for section"`
	OnSuccess string `long:"on-success" description:"optional message for successful subcommand"`
	OnFailure string `long:"on-failure" description:"optional message for failed subcommand"`
}

type SectionOpt struct {
	Section
	Out   io.Writer
	Clock clock.Clock
	Exec  exec.Exec
}

type SectionError struct {
	Retval int
	Err    error
}

func writeSection(opts SectionOpt) error {
	machineLog := &mrl.MachineReadableLog{
		Name:   opts.Name,
		Type:   fmt.Sprintf("section-%s", opts.Type),
		Result: opts.Result,
		Time:   opts.Clock.Now(),
	}

	var humanReadable string
	if opts.Type == "start" {
		humanReadable = fmt.Sprintf("section-%s: '%s'",
			opts.Type,
			opts.Name)
	} else if opts.Type == "end" {
		if opts.Result == 0 && opts.OnSuccess != "" {
			humanReadable = fmt.Sprintf("section-%s: '%s' result: %d message: '%s'",
				opts.Type,
				opts.Name,
				opts.Result,
				opts.OnSuccess)
			machineLog.Message = opts.OnSuccess
		} else if opts.Result != 0 && opts.OnFailure != "" {
			humanReadable = fmt.Sprintf("section-%s: '%s' result: %d message: '%s'",
				opts.Type,
				opts.Name,
				opts.Result,
				opts.OnFailure)
			machineLog.Message = opts.OnFailure
		} else {
			humanReadable = fmt.Sprintf("section-%s: '%s' result: %d",
				opts.Type,
				opts.Name,
				opts.Result)
		}
	} else {
		return errors.New("invalid section type argument")
	}

	_, err := fmt.Fprint(opts.Out, humanReadable)
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	machineLogJSON, err := json.Marshal(machineLog)
	if err != nil { // !branch-not-tested
		return err
	}

	_, err = fmt.Fprintf(opts.Out, " MRL:%s\n", string(machineLogJSON))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func (e *SectionError) Unwrap() error { return e.Err }
func (e *SectionError) Error() string {
	// This returns an empty string so it does not print the error outside
	// of the section-end block. The error information is printed inside the
	// SectionOpt.Execute function below
	return ""
}

func (opts *SectionOpt) Execute(args []string) error {
	if opts.Name == "" {
		return errors.New("missing section name")
	}

	if opts.Type == "section" {
		if len(args) == 0 {
			return errors.New("the section subcommand requires a command parameter '-- <command> ...'")
		}

		sectionOpts := *opts
		sectionOpts.Type = "start"
		if err := writeSection(sectionOpts); err != nil {
			return err
		}

		cmd := opts.Exec.Command(args[0], args[1:]...)
		cmd.SetOutput(opts.Out)
		err := cmd.Run()

		exitCode := 0

		var sectionError *SectionError

		if err != nil {
			var e *os_exec.ExitError
			if errors.As(err, &e) {
				exitCode = e.ExitCode()
			} else {
				exitCode = -1
			}
			fmt.Fprintf(opts.Out, "Section subcommand failed with %d: %s\n", exitCode, err)
			sectionError = &SectionError{exitCode, err}
		}
		sectionOpts.Type = "end"
		sectionOpts.Result = exitCode
		err = writeSection(sectionOpts)

		if sectionError != nil {
			// returning sectionError to propagate the resulting non-zero result code
			return sectionError
		}
		return err
	}

	return writeSection(*opts)
}
