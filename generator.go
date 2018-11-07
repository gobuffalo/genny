package genny

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/gobuffalo/events"
	"github.com/gobuffalo/packd"
	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generator is the basic type for generators to use
type Generator struct {
	StepName     string
	Should       func(*Runner) bool
	Root         string
	ErrorFn      func(error)
	runners      []RunFn
	transformers []Transformer
	moot         *sync.RWMutex
}

func stepName() string {
	bb := &bytes.Buffer{}
	for i := 0; i < 5; i++ {
		_, file, line, _ := runtime.Caller(i)
		mod := time.Now()
		if info, err := os.Stat(file); err == nil {
			mod = info.ModTime()
		}
		bb.WriteString(fmt.Sprintf("%s:%d:%d\n", file, line, mod.UnixNano()))
	}
	h := sha1.New()
	h.Write(bb.Bytes())
	return fmt.Sprintf("%x", h.Sum(nil))[:8]
}

// New, well-formed, generator
func New() *Generator {
	g := &Generator{
		StepName:     stepName(),
		runners:      []RunFn{},
		moot:         &sync.RWMutex{},
		transformers: []Transformer{},
	}
	return g
}

func (g *Generator) Event(kind string, payload events.Payload) {
	g.RunFn(func(r *Runner) error {
		return events.EmitPayload(kind, payload)
	})
}

// File adds a file to be run when the generator is run
func (g *Generator) File(f File) {
	g.RunFn(func(r *Runner) error {
		return r.File(f)
	})
}

func (g *Generator) Transform(f File) (File, error) {
	g.moot.RLock()
	defer g.moot.RUnlock()
	var err error
	for _, t := range g.transformers {
		f, err = t.Transform(f)
		if err != nil {
			return f, errors.WithStack(err)
		}
	}

	return f, nil
}

// Transformer adds a file transform to the generator
func (g *Generator) Transformer(t Transformer) {
	g.moot.Lock()
	defer g.moot.Unlock()
	g.transformers = append(g.transformers, t)
}

// Command adds a command to be run when the generator is run
func (g *Generator) Command(cmd *exec.Cmd) {
	g.RunFn(func(r *Runner) error {
		return r.Exec(cmd)
	})
}

// Box walks through a packr.Box and adds Files for each entry
// in the box.
func (g *Generator) Box(box packd.Walker) error {
	return box.Walk(func(path string, f packd.File) error {
		g.File(NewFile(path, f))
		return nil
	})
}

// RunFn adds a generic "runner" function to the generator.
func (g *Generator) RunFn(fn RunFn) {
	g.moot.Lock()
	defer g.moot.Unlock()
	g.runners = append(g.runners, fn)
}

func (g1 *Generator) Merge(g2 *Generator) {
	g2.moot.Lock()
	g1.moot.Lock()
	g1.runners = append(g1.runners, g2.runners...)
	g1.transformers = append(g1.transformers, g2.transformers...)
	g1.moot.Unlock()
	g2.moot.Unlock()
}
