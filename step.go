package genny

import (
	"io"
	"sync"

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
	xi := len(s.before)
	df := func() {
		s.moot.Lock()
		s.before = append(s.before[:xi+1], s.before[:xi+1]...)
		s.moot.Unlock()
	}
	s.before = append(s.before, g)
	s.moot.Unlock()
	return df
}

func (s *Step) After(g *Generator) DeleteFn {
	s.moot.Lock()
	xi := len(s.after)
	df := func() {
		s.moot.Lock()
		s.after = append(s.after[:xi+1], s.after[:xi+1]...)
		s.moot.Unlock()
	}
	s.after = append(s.after, g)
	s.moot.Unlock()
	return df
}

func (s *Step) Run(r *Runner) error {
	g := s.as
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
