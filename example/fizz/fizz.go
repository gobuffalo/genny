package fizz

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/markbates/inflect"
)

func FizzMigration(ats attrs.Attrs) (*genny.Generator, error) {
	g := genny.New()
	if len(ats) < 1 {
		return g, errors.New("requires at least attr")
	}

	if err := g.Box(packr.NewBox("./templates")); err != nil {
		return g, err
	}

	data := struct {
		Model inflect.Name
		Attrs attrs.Attrs
	}{
		Model: ats[0].Name,
	}
	if len(ats) > 1 {
		data.Attrs = ats[1:]
	}

	g.Transformer(gotools.TemplateTransformer(data, template.FuncMap{}))
	g.Transformer(genny.NewTransformer(".fizz", func(f genny.File) (genny.File, error) {
		if f.Name() != "migration.fizz" {
			return f, nil
		}
		t := time.Now()
		return genny.NewFile(filepath.Join("migrations", fmt.Sprintf("%d_create_%s.fizz", t.UnixNano(), data.Model.PluralUnder())), f), nil
	}))
	return g, nil
}
