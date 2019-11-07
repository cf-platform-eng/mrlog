package mrlog

import (
	os_exec "os/exec"
	"io"
	"github.com/cf-platform-eng/mrlog/exec"
)

type Cmd struct {
	Cmd *os_exec.Cmd
}

func (cmd *Cmd) SetOutput(writer io.Writer) {
	cmd.Cmd.Stdout = writer
	cmd.Cmd.Stderr = writer
}

func (cmd *Cmd) Run() error {
	return cmd.Cmd.Run()
}

type Exec struct {
}

func (e *Exec) Command(command string, arg ...string) exec.Cmd {
	newCmd := Cmd{}
	newCmd.Cmd = os_exec.Command(command, arg...)
	return &newCmd
}