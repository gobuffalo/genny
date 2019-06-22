package gogen

import (
	"bytes"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen/goimports"
)

func Fmt(root string) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		i, err := goimports.New(root)
		if err != nil {
			return err
		}
		return i.Run()
	})

	return g, nil
}

func FmtTransformer() genny.Transformer {
	t := genny.NewTransformer(".go", func(f genny.File) (genny.File, error) {
		bb := &bytes.Buffer{}
		gi := goimports.NewFromFiles(goimports.File{
			Name: f.Name(),
			In:   f,
			Out:  bb,
		})
		if err := gi.Run(); err != nil {
			return f, err
		}
		f = genny.NewFile(f.Name(), bb)
		return f, nil
	})
	t.StripExt = false
	return t
}
