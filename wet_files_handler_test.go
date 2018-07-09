package genny

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (r *Suite) Test_WetFilesHandler() {
	g := Background()

	f1 := filepath.Join(r.PWD, "foo.txt")
	c1 := strings.ToUpper(f1)
	f2 := filepath.Join(r.PWD, "bar/baz.txt")
	c2 := strings.ToUpper(f2)

	g = WithFileFromReader(g, f1, strings.NewReader(c1))
	g = WithFileFromReader(g, f2, strings.NewReader(c2))
	g = WetFilesHandler(g)
	r.NoError(g.Run())

	b, err := ioutil.ReadFile(f1)
	r.NoError(err)
	r.Equal(c1, string(b))

	b, err = ioutil.ReadFile(f2)
	r.NoError(err)
	r.Equal(c2, string(b))
}
