package exec

import (
	"io"
)

//go:generate counterfeiter Exec
type Exec interface {
	Command(command string, arg ...string) Cmd
}

//go:generate counterfeiter Cmd
type Cmd interface {
	SetOutput(writer io.Writer)
	Run() error
}