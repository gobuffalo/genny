package genny

import (
	"io"
)

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
