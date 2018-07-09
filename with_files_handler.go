package genny

import (
	"github.com/pkg/errors"
)

// FileHandlerFn used by `WithFilesHandlers`
type FileHandlerFn func(f File) error

type withFilesHandler struct {
	Generator
	fn FileHandlerFn
}

func (w withFilesHandler) Parent() Generator {
	return w.Generator
}

func (w withFilesHandler) String() string {
	return "genny.WithFilesHandler"
}

func (w withFilesHandler) Run() error {
	files := Files(w.Generator)
	for _, f := range files {
		if err := w.fn(f); err != nil {
			return errors.WithStack(err)
		}
	}
	return w.Parent().Run()
}

// WithFilesHandler processes all of the files up the generator
// tree with the same handler function
func WithFilesHandler(g Generator, fn FileHandlerFn) Generator {
	g = withFilesHandler{
		Generator: g,
		fn:        fn,
	}
	return g
}
