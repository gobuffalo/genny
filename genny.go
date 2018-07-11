package genny

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// New returns a basic, "background", generator
// with the given context. If you want a generator
// with a context.Background use `Background()`
// instead.
func New(ctx context.Context) Generator {
	g := WithContext(background{}, ctx)

	l := logrus.New()
	g = WithLogger(g, l)
	return g
}

func Run(g Generator) error {
	gens := Tree(g)
	rgens := make(Generators, len(gens))
	for i, x := range gens {
		rgens[len(gens)-i-1] = x
	}
	for _, g := range rgens {
		if err := g.Run(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
