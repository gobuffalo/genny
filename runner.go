package genny

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Runner will run the generators
type Runner struct {
	Logger     Logger                // Logger to use for the run
	Context    context.Context       // context to use for the run
	ExecFn     func(*exec.Cmd) error // function to use when executing files
	FileFn     func(File) error      // function to use when writing files
	Root       string                // the root of the write path
	generators []*Generator
	moot       *sync.Mutex
	results    Results
}

func (r *Runner) Results() Results {
	return r.results
}

func (r *Runner) WithRun(fn RunFn) {
	g := New()
	g.RunFn(fn)
	r.With(g)
}

// With adds a Generator to the Runner
func (r *Runner) With(g *Generator) {
	r.moot.Lock()
	defer r.moot.Unlock()
	r.generators = append(r.generators, g)
}

// WithFn will evaluate the function and if successful it will add
// the Generator to the Runner, otherwise it will return the error
func (r *Runner) WithFn(fn func() (*Generator, error)) error {
	g, err := fn()
	if err != nil {
		return errors.WithStack(err)
	}
	r.With(g)
	return nil
}

// Run all of the generators!
func (r *Runner) Run() error {
	r.moot.Lock()
	defer r.moot.Unlock()
	for _, g := range r.generators {
		for _, fn := range g.runners {
			if err := fn(r); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

// Exec can be used inside of Generators to run commands
func (r *Runner) Exec(cmd *exec.Cmd) error {
	r.results.Commands = append(r.results.Commands, cmd)
	r.Logger.Infof(strings.Join(cmd.Args, " "))
	if r.ExecFn == nil {
		return nil
	}
	return r.ExecFn(cmd)
}

// File can be used inside of Generators to write files
func (r *Runner) File(f File) error {
	r.results.Files = append(r.results.Files, f)
	name := f.Name()
	if !filepath.IsAbs(name) {
		name = filepath.Join(r.Root, name)
	}
	r.Logger.Infof(name)
	if r.FileFn == nil {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}
		r.Logger.Debugf(string(b))
		return nil
	}
	return r.FileFn(f)
}

type RunFn func(r *Runner) error
