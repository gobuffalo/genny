package genny

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type RunFn func(r *Runner) error

// Runner will run the generators
type Runner struct {
	Logger     Logger                           // Logger to use for the run
	Context    context.Context                  // context to use for the run
	ExecFn     func(*exec.Cmd) error            // function to use when executing files
	FileFn     func(File) (File, error)         // function to use when writing files
	ChdirFn    func(string, func() error) error // function to use when changing directories
	DeleteFn   func(string) error               // function used to delete files/folders
	Root       string                           // the root of the write path
	Disk       *Disk
	generators []*Generator
	moot       *sync.RWMutex
	results    Results
}

func (r *Runner) Results() Results {
	r.moot.Lock()
	defer r.moot.Unlock()
	r.results.Files = r.Disk.Files()
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
		if g.Should != nil {
			if !g.Should(r) {
				continue
			}
		}
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
	name := f.Name()
	if !filepath.IsAbs(name) {
		name = filepath.Join(r.Root, name)
	}
	r.Logger.Infof(name)
	if r.FileFn != nil {
		var err error
		if f, err = r.FileFn(f); err != nil {
			return errors.WithStack(err)
		}
		if s, ok := f.(io.Seeker); ok {
			s.Seek(0, 0)
		}
	}
	f = NewFile(f.Name(), f)
	if s, ok := f.(io.Seeker); ok {
		s.Seek(0, 0)
	}
	r.Disk.Add(f)
	return nil
}

func (r *Runner) FindFile(name string) (File, error) {
	return r.Disk.Find(name)
}

// Chdir will change to the specified directory
// and revert back to the current directory when
// the runner function has returned.
// If the directory does not exist, it will be
// created for you.
func (r *Runner) Chdir(path string, fn func() error) error {
	r.Logger.Infof("CD: %s", path)

	if r.ChdirFn != nil {
		return r.ChdirFn(path, fn)
	}

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	os.MkdirAll(path, 0755)
	if err := os.Chdir(path); err != nil {
		return errors.WithStack(err)
	}
	if err := fn(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *Runner) Delete(path string) error {
	r.Logger.Infof("rm: %s", path)

	defer r.Disk.Remove(path)
	if r.DeleteFn != nil {
		return r.DeleteFn(path)
	}
	return nil
}
