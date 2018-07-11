package genny

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func (r *Suite) Test_WithFileTransformer() {
	g := Background()

	g = WithFileTransformer(g, func(f File) (File, error) {
		f = NewFile(f.Name(), strings.NewReader("transformed!"))
		return f, nil
	})

	var finished string
	g = WithFilesHandler(g, func(f File) error {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}
		finished = string(b)
		return nil
	})

	g = WithFileFromReader(g, "foo.txt", strings.NewReader("asdfasfd"))

	r.NoError(Run(g))
	r.Equal("transformed!", finished)
}
