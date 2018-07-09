package gotools

import (
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_WithTemplate(t *testing.T) {
	r := require.New(t)

	input := `hello {{.}}`
	g := genny.Background()
	g = genny.WithFile(g, genny.NewFile("foo.go", strings.NewReader(input)))

	g, err := WithTemplate(g, "mark", template.FuncMap{})
	r.NoError(err)

	fa, ok := g.(genny.Fileable)
	r.True(ok)
	b, err := ioutil.ReadAll(fa.File())
	r.NoError(err)

	r.Equal("hello mark", string(b))
}

func Test_WithTemplate_NoFile(t *testing.T) {
	r := require.New(t)

	g := genny.Background()

	_, err := WithTemplate(g, "mark", template.FuncMap{})
	r.Error(err)
	r.Equal(genny.ErrNilFile, err)
}

func Test_WithTemplate_BadTemplate(t *testing.T) {
	r := require.New(t)

	input := `hello {{.}`
	g := genny.Background()
	g = genny.WithFile(g, genny.NewFile("foo.go", strings.NewReader(input)))

	_, err := WithTemplate(g, "mark", template.FuncMap{})
	r.Error(err)
}
