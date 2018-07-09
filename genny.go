package genny

import (
	"context"

	"github.com/sirupsen/logrus"
)

// New returns a basic, "background", generator
// with the given context. If you want a generator
// with a context.Background use `Background()`
// instead.
func New(ctx context.Context) Generator {
	g := Context(background{}, ctx)

	l := logrus.New()
	g = &withLogger{
		Generator: g,
		logger:    l,
	}

	return g
}
