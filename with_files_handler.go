package genny

import (
	"github.com/pkg/errors"
)

type FileHandler interface {
	Handle(File) error
}

// FileHandlerFn used by `WithFilesHandlers`
type FileHandlerFn func(f File) error

type withFilesHandler struct {
	Generator
	fn FileHandlerFn
}

func (w withFilesHandler) Handle(f File) error {
	w.Logger().Infof("[genny.WithFilesHandler] %s", f.Name())
	var err error
	for _, t := range Tree(w) {
		if ft, ok := t.(FileTransformer); ok {
			f, err = ft.Transform(f)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	if w.fn == nil {
		return nil
	}
	return w.fn(f)
}

func (w withFilesHandler) Parent() Generator {
	return w.Generator
}

func (w withFilesHandler) String() string {
	return "genny.WithFilesHandler"
}

// WithFilesHandler processes all of the files with the handler function
func WithFilesHandler(g Generator, fn FileHandlerFn) Generator {
	g = withFilesHandler{
		Generator: g,
		fn:        fn,
	}
	return g
}
