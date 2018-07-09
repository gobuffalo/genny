package plushgen

import (
	"io/ioutil"
	"strings"
	"testing"

	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/require"
)

func Test_WithTemplate(t *testing.T) {
	r := require.New(t)

	g := genny.Background()
	g = genny.WithFile(g, genny.NewFile("foo.go", strings.NewReader(`{{.}} <%= name %>`)))

	g, err := gotools.WithTemplate(g, "Hello", template.FuncMap{})
	r.NoError(err)

	ctx := plush.NewContext()
	ctx.Set("name", "mark")
	g, err = WithTemplate(g, ctx)
	r.NoError(err)

	fa, ok := g.(genny.Fileable)
	r.True(ok)
	b, err := ioutil.ReadAll(fa.File())
	r.NoError(err)

	r.Equal("Hello mark", string(b))
}
