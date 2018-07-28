package pop

import (
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

func Model(ats attrs.NamedAttrs) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(packr.NewBox("./templates")); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	ctx.Set("model", ats)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.NewTransformer(".go", func(f genny.File) (genny.File, error) {
		return genny.NewFile(filepath.Join("models", ats.Name.File(".go").String()), f), nil
	}))
	return g, nil
}
