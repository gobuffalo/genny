package genny

import (
	"os/exec"
	"strings"
)

type Execer interface {
	Exec(*exec.Cmd) error
}

type withExec struct {
	Generator
	fn ExecFn
}

func (d withExec) Parent() Generator {
	return d.Generator
}

func (d withExec) Exec(cmd *exec.Cmd) error {
	d.Logger().Infof("[genny.WithExec] %s", strings.Join(cmd.Args, " "))
	if d.fn == nil {
		return nil
	}
	return d.fn(cmd)
}

func (withExec) String() string {
	return "genny.WithExec"
}

// ExecFn used with WithExec
type ExecFn func(*exec.Cmd) error

// WithExec will run commands. This should be placed at the BEGINNING
// of your generator stack.
func WithExec(g Generator, fn ExecFn) Generator {
	return withExec{
		Generator: g,
		fn:        fn,
	}
}
