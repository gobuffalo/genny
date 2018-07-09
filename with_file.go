package genny

import "io"

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
