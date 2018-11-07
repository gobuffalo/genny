package genny

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gobuffalo/events"
	"github.com/markbates/safe"
	"github.com/pkg/errors"
)

type DeleteFn func()

type Step struct {
	as     *Generator
	before []*Generator
	after  []*Generator
	index  int
	moot   *sync.RWMutex
}

func (s *Step) Before(g *Generator) DeleteFn {
	s.moot.Lock()
	df := func() {
		var a []*Generator
		s.moot.Lock()
		for _, b := range s.before {
			if g.StepName == b.StepName {
				continue
			}
			a = append(a, b)
		}
		s.before = a
		s.moot.Unlock()
	}
	s.before = append(s.before, g)
	s.moot.Unlock()
	return df
}

func (s *Step) After(g *Generator) DeleteFn {
	s.moot.Lock()
	df := func() {
		var a []*Generator
		s.moot.Lock()
		for _, b := range s.after {
			if g.StepName == b.StepName {
				continue
			}
			a = append(a, b)
		}
		s.after = a
		s.moot.Unlock()
	}
	s.after = append(s.after, g)
	s.moot.Unlock()
	return df
}

func (s *Step) Run(r *Runner) error {
	for _, b := range s.before {
		if err := s.runGenerator(r, b); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := s.runGenerator(r, s.as); err != nil {
		return errors.WithStack(err)
	}

	for _, b := range s.after {
		if err := s.runGenerator(r, b); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Step) runGenerator(r *Runner, g *Generator) error {
	r.curGen = g

	payload := events.Payload{
		"runner":    r,
		"step":      s,
		"generator": g,
	}
	if g.Should != nil {
		err := safe.RunE(func() error {
			if !g.Should(r) {
				return io.EOF
			}
			return nil
		})
		if err != nil {
			r.Logger.Debugf("Step: %s [skipped]", g.StepName)
			events.EmitPayload(EvtStepPrefix+":skipping:"+g.StepName, payload)
			return nil
		}
	}
	r.Logger.Debugf("Step: %s", g.StepName)
	events.EmitPayload(EvtStepPrefix+":running:"+g.StepName, payload)
	return r.Chdir(r.Root, func() error {
		for _, fn := range g.runners {
			err := safe.RunE(func() error {
				return fn(r)
			})
			if err != nil {
				events.EmitError(EvtStepPrefix+":running:"+g.StepName+":err", err, payload)
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

func NewStep(g *Generator, index int) (*Step, error) {
	if g == nil {
		return nil, errors.New("generator can not be nil")
	}
	return &Step{
		as:    g,
		index: index,
		moot:  &sync.RWMutex{},
	}, nil
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
	bb.WriteString(fmt.Sprint(time.Now().UnixNano()))
	h := sha1.New()
	h.Write(bb.Bytes())
	return fmt.Sprintf("%x", h.Sum(nil))[:8]
}
