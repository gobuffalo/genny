package plushgen

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/require"
)

func Test_Transformer(t *testing.T) {
	r := require.New(t)

	ctx := plush.NewContext()
	ctx.Set("name", "mark")
	f := genny.NewFile("foo.plush.txt", strings.NewReader("Hello <%= name %>"))

	tr := Transformer(ctx)
	f, err := tr.Transform(f)
	r.NoError(err)
	r.Equal("foo.txt", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("Hello mark", string(b))
}

func Test_Transformer_No_Ext(t *testing.T) {
	r := require.New(t)

	ctx := plush.NewContext()
	ctx.Set("name", "mark")
	f := genny.NewFile("foo.txt", strings.NewReader("Hello <%= name %>"))

	tr := Transformer(ctx)
	f, err := tr.Transform(f)
	r.NoError(err)
	r.Equal("foo.txt", f.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("Hello <%= name %>", string(b))
}
