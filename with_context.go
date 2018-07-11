package genny

import (
	"context"
)

type withContext struct {
	Generator
	ctx context.Context
}

func (w withContext) Context() context.Context {
	return w.ctx
}

func (w withContext) Parent() Generator {
	return w.Generator
}

func (w withContext) String() string {
	return "genny.WithContext"
}

// Context returns a new Generator with the given context
func WithContext(g Generator, ctx context.Context) Generator {
	return withContext{
		Generator: g,
		ctx:       ctx,
	}
}
