package plushgen

import (
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

// WithTemplate returns a generator who's previous file has been run through plush
func WithTemplate(g genny.Generator, ctx *plush.Context) (genny.Generator, error) {
	fa, ok := g.(genny.Fileable)
	if !ok {
		return g, genny.ErrNilFile
	}
	f, err := renderWithPlush(fa.File(), ctx)
	if err != nil {
		return g, errors.WithStack(err)
	}
	return genny.WithFile(g, f), nil
}

func renderWithPlush(f genny.File, ctx *plush.Context) (genny.File, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return f, errors.WithStack(err)
	}
	s, err := plush.Render(string(b), ctx)
	if err != nil {
		return f, errors.WithStack(err)
	}
	return genny.NewFile(f.Name(), strings.NewReader(s)), nil
}
