package plushgen

import (
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

// WithPlush adds a file transform to render files with plush syntax
func WithPlush(g genny.Generator, ctx *plush.Context) genny.Generator {
	g = genny.WithFileTransformer(g, func(f genny.File) (genny.File, error) {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return f, errors.WithStack(err)
		}
		s, err := plush.Render(string(b), ctx)
		if err != nil {
			return f, errors.WithStack(err)
		}
		return genny.NewFile(f.Name(), strings.NewReader(s)), nil
	})
	return g
}
