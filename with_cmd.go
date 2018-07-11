package genny

import (
	"os/exec"

	"github.com/pkg/errors"
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

func (w withCmd) Run() error {
	// if err := w.Parent().Run(); err != nil {
	// 	return errors.WithStack(err)
	// }
	for _, p := range Tree(w) {
		if e, ok := p.(Execer); ok {
			err := e.Exec(w.Cmd())
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
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
