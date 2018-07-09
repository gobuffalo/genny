package genny

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// ExecFn used with WithExec
type ExecFn func(*exec.Cmd) error

type withExec struct {
	Generator
	exec ExecFn
}

func (w withExec) Parent() Generator {
	return w.Generator
}

func (w withExec) String() string {
	return "genny.WithExec"
}

func (w withExec) Run() error {
	for _, c := range Cmds(w.Generator) {
		w.Logger().Infof("[%s] %s", w.String(), strings.Join(c.Args, " "))
		if err := w.exec(c); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// WithExec processes all of the files up the generator
// tree with the same handler function
func WithExec(g Generator, fn ExecFn) Generator {
	g = withExec{
		Generator: g,
		exec:      fn,
	}
	return g
}
