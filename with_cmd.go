package genny

import (
	"os/exec"
)

type withCmd struct {
	Generator
	cmd *exec.Cmd
}

func (w withCmd) Parent() Generator {
	return w.Generator
}

func (w withCmd) Cmd() *exec.Cmd {
	return w.cmd
}

func (w withCmd) String() string {
	return "genny.WithCmd"
}

// WithCmd wraps the generator with a command to be executed
// during the run of the generator.
func WithCmd(g Generator, cmd *exec.Cmd) Generator {
	g = withCmd{
		Generator: g,
		cmd:       cmd,
	}
	return g
}
