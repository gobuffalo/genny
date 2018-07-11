package gotools

import (
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_WithTemplate(t *testing.T) {
	r := require.New(t)

	g := genny.Background()

	g = WithTemplate(g, "Hello", template.FuncMap{})

	var finished string
	g = genny.WithFilesHandler(g, func(f genny.File) error {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}
		finished = string(b)
		return nil
	})

	g = genny.WithFile(g, genny.NewFile("foo.go", strings.NewReader(`{{.}} mark`)))

	r.NoError(genny.Run(g))

	r.Equal("Hello mark", finished)
}
