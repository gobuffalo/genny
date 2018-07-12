package pop

import (
	"errors"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

func Model(ats attrs.Attrs) (*genny.Generator, error) {
	g := genny.New()
	if len(ats) < 1 {
		return g, errors.New("requires at least attr")
	}

	if err := g.Box(packr.NewBox("./templates")); err != nil {
		return g, err
	}

	name := ats[0].Name
	ctx := plush.NewContext()
	ctx.Set("model", ats[0])
	if len(ats) > 1 {
		ctx.Set("attrs", ats[1:])
	}

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.NewTransformer(".go", func(f genny.File) (genny.File, error) {
		return genny.NewFile(filepath.Join("models", name.File()+".go"), f), nil
	}))
	return g, nil
}
