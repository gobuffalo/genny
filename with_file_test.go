package genny

import (
	"io/ioutil"
	"strings"
)

func (r *Suite) Test_WithFile() {
	g := Background()
	g = WithFile(g, NewFile("foo.txt", strings.NewReader("foo")))

	fa, ok := g.(Fileable)
	r.True(ok)

	f := fa.File()
	r.NotNil(f)
	r.Equal("foo.txt", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("foo", string(b))
}

func (r *Suite) Test_WithFileFromReader() {
	g := Background()
	g = WithFileFromReader(g, "foo.txt", strings.NewReader("foo"))

	fa, ok := g.(Fileable)
	r.True(ok)

	f := fa.File()
	r.NotNil(f)
	r.Equal("foo.txt", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("foo", string(b))
}
