package pop

import (
	"errors"
	"path/filepath"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

func Model(attrs attrs.Attrs) (*genny.Generator, error) {
	g := genny.New()
	if len(attrs) < 1 {
		return g, errors.New("requires at least attr")
	}

	if err := g.Box(packr.NewBox("./templates")); err != nil {
		return g, err
	}

	name := attrs[0].Name
	ctx := plush.NewContext()
	ctx.Set("model", attrs[0])
	if len(attrs) > 1 {
		ctx.Set("attrs", attrs[1:])
	}

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.NewTransformer(".go", func(f genny.File) (genny.File, error) {
		return genny.NewFile(filepath.Join("models", name.File()+".go"), f), nil
	}))
	return g, nil
}
