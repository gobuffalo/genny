package genny

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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
	FileFn     func(File) error                 // function to use when writing files
	ChdirFn    func(string, func() error) error // function to use when changing directories
	Root       string                           // the root of the write path
	generators []*Generator
	moot       *sync.RWMutex
	results    Results
	files      map[string]File
}

func (r *Runner) Results() Results {
	r.moot.Lock()
	defer r.moot.Unlock()
	var files []File
	for _, f := range r.files {
		files = append(files, f)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	r.results.Files = files
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
	defer func() {
		r.files[f.Name()] = f
	}()
	name := f.Name()
	if !filepath.IsAbs(name) {
		name = filepath.Join(r.Root, name)
	}
	r.Logger.Infof(name)
	if r.FileFn != nil {
		return r.FileFn(f)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.WithStack(err)
	}
	f = NewFile(f.Name(), bytes.NewReader(b))
	r.Logger.Debugf(string(b))
	return nil
}

func (r *Runner) FindFile(name string) (File, error) {
	if f, ok := r.files[name]; ok {
		if seek, ok := f.(io.Seeker); ok {
			seek.Seek(0, 0)
		}
		return f, nil
	}

	gf := NewFile(name, bytes.NewReader([]byte("")))
	f, err := os.Open(name)
	if err != nil {
		return gf, errors.WithStack(err)
	}
	defer f.Close()

	bb := &bytes.Buffer{}

	if _, err := io.Copy(bb, f); err != nil {
		return gf, errors.WithStack(err)
	}

	return NewFile(name, bb), nil
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
