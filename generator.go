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
