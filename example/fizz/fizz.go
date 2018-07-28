package fizz

import (
	"fmt"
	"path/filepath"
	"time"

	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
)

func FizzMigration(ats attrs.NamedAttrs) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(packr.NewBox("./templates")); err != nil {
		return g, err
	}

	g.Transformer(gotools.TemplateTransformer(ats, template.FuncMap{}))
	g.Transformer(genny.NewTransformer(".fizz", func(f genny.File) (genny.File, error) {
		if f.Name() != "migration.fizz" {
			return f, nil
		}
		t := time.Now()
		p := ats.Name.File().Pluralize().String()
		return genny.NewFile(filepath.Join("migrations", fmt.Sprintf("%d_create_%s.fizz", t.UnixNano(), p)), f), nil
	}))
	return g, nil
}
