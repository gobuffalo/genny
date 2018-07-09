package genny

import (
	"io"

	"github.com/pkg/errors"
)

// ErrNilFile should be returned when it is
// expected that the parent generator is supposed
// to have a file.
var ErrNilFile = errors.New(".File() is nil")

// File interface for working with files
type File interface {
	io.Reader
	Name() string
}

type simpleFile struct {
	io.Reader
	name string
}

func (s simpleFile) Name() string {
	return s.name
}

// NewFile takes the name of the file you want to
// write to and a reader to reader from
func NewFile(name string, r io.Reader) File {
	return simpleFile{
		Reader: r,
		name:   name,
	}
}

// Files returns all of the files up to the parent
// generator.
func Files(g Generator) []File {
	var files []File
	t := map[string]bool{}
	fp := func(f File) {
		if f != nil {
			if _, ok := t[f.Name()]; !ok {
				files = append([]File{f}, files...)
				t[f.Name()] = true
			}
		}
	}
	if fa, ok := g.(Fileable); ok {
		fp(fa.File())
	}

	p := g.Parent()
	for p != nil {
		if fa, ok := p.(Fileable); ok {
			fp(fa.File())
		}
		p = p.Parent()
	}
	return files
}

// Fileable interface for generators that return a file
type Fileable interface {
	File() File
}
