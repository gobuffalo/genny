package new

import (
	"path"
	"path/filepath"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("genny:genny:new", "../new/templates")); err != nil {
		return g, errors.WithStack(err)
	}
	name := name.New(opts.Name)
	ctx := plush.NewContext()
	ctx.Set("name", name)
	pkg, err := envy.CurrentModule()
	if err != nil {
		return g, errors.WithStack(err)
	}
	ctx.Set("boxName", path.Join(pkg, opts.Prefix, opts.Name))
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("-name-", name.File().String()))
	g.Transformer(genny.Dot())
	g.Transformer(genny.NewTransformer("*", func(f genny.File) (genny.File, error) {
		f = genny.NewFile(filepath.Join(opts.Prefix, f.Name()), f)
		return f, nil
	}))
	return g, nil
}
