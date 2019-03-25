package command

import (
	"path"

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

	if err := g.Box(packr.New("github.com/gobuffalo/genny/genny/command/templates", "../command/templates")); err != nil {
		return g, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)

	pkg := path.Join(opts.App.PackagePkg, opts.Name)
	ctx.Set("genPkg", path.Join(pkg, opts.Prefix, opts.Name))
	ctx.Set("cmdPkg", path.Join(pkg, "cmd"))
	g.Transformer(plushgen.Transformer(ctx))
	return g, nil
}
