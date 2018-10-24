package genny

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
)

// File interface for working with files
type File interface {
	fmt.Stringer
	io.Reader
	Name() string
	io.Writer
}

var _ io.ReadWriter = &simpleFile{}

type simpleFile struct {
	io.Reader
	name string
}

func (s simpleFile) Name() string {
	return s.name
}

func (s *simpleFile) Write(p []byte) (int, error) {
	bb := &bytes.Buffer{}
	i, err := bb.Write(p)
	s.Reader = bb
	return i, err
}

func (s *simpleFile) String() string {
	src, _ := ioutil.ReadAll(s)
	s.Reader = bytes.NewReader(src)
	return string(src)
}

func (s simpleFile) Seek(offset int64, whence int) (int64, error) {
	if seek, ok := s.Reader.(io.Seeker); ok {
		return seek.Seek(offset, whence)
	}
	return -1, nil
}

// NewFile takes the name of the file you want to
// write to and a reader to reader from
func NewFile(name string, r io.Reader) File {
	osname := name
	if runtime.GOOS == "windows" {
		osname = strings.Replace(osname, "\\", "/", -1)
	}

	if r == nil {
		r = &bytes.Buffer{}
	}
	if seek, ok := r.(io.Seeker); ok {
		seek.Seek(0, 0)
	}
	return &simpleFile{
		Reader: r,
		name:   osname,
	}
}

func NewFileS(name string, s string) File {
	return NewFile(name, strings.NewReader(s))
}

func NewFileB(name string, s []byte) File {
	return NewFile(name, bytes.NewReader(s))
}
