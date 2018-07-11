package genny

import (
	"context"
)

// Generator is the main interface for writing new generators
type Generator interface {
	Context() context.Context
	Run() error
	Parent() Generator
	Logger() Logger
}

type Generators []Generator

func Tree(g Generator) Generators {
	gens := Generators{g}
	p := g.Parent()
	for p != nil {
		gens = append(gens, p)
		p = p.Parent()
	}
	return gens
}
