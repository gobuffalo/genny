package genny

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewStep(t *testing.T) {
	r := require.New(t)

	_, err := NewStep(nil, 0)
	r.Error(err)

	s, err := NewStep(New(), 0)
	r.NoError(err)
	r.NotZero(s)
}

func Test_StepBefore(t *testing.T) {
	r := require.New(t)

	var actual []string

	g := New()
	g.RunFn(func(r *Runner) error {
		actual = append(actual, "as")
		return nil
	})

	s, err := NewStep(g, 0)
	r.NoError(err)

	b := New()
	b.RunFn(func(r *Runner) error {
		actual = append(actual, "before")
		return nil
	})
	s.Before(b)

	run := DryRunner(context.Background())

	r.NoError(s.Run(run))

	r.Equal([]string{"before", "as"}, actual)
}

func Test_StepAfter(t *testing.T) {
	r := require.New(t)

	var actual []string

	g := New()
	g.RunFn(func(r *Runner) error {
		actual = append(actual, "as")
		return nil
	})

	s, err := NewStep(g, 0)
	r.NoError(err)

	b := New()
	b.RunFn(func(r *Runner) error {
		actual = append(actual, "after")
		return nil
	})
	s.After(b)

	run := DryRunner(context.Background())

	r.NoError(s.Run(run))

	r.Equal([]string{"as", "after"}, actual)
}
