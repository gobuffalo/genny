package genny

import (
	"io"

	"github.com/pkg/errors"
)

type withFile struct {
	Generator
	file File
}

func (w withFile) Parent() Generator {
	return w.Generator
}

func (w withFile) File() File {
	return w.file
}

func (w withFile) String() string {
	return "genny.WithFile"
}

func (w withFile) Run() error {
	for _, p := range Tree(w) {
		if e, ok := p.(FileHandler); ok {
			err := e.Handle(w.File())
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

// WithFile returns a generator wrapped with a file
func WithFile(g Generator, f File) Generator {
	return withFile{
		Generator: g,
		file:      f,
	}
}

// WithFileFromReader returns a generator wrapped with a file from the reader
func WithFileFromReader(g Generator, name string, r io.Reader) Generator {
	return withFile{
		Generator: g,
		file:      NewFile(name, r),
	}
}
