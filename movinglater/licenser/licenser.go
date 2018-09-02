package licenser

import (
	"bytes"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

var Available []string
var box packr.Box

func init() {
	box = packr.NewBox("../licenser/templates")
	box.Walk(func(path string, f packr.File) error {
		name := filepath.Base(path)
		Available = append(Available, name)
		return nil
	})
}

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	body, err := box.MustBytes(opts.Name)
	if err != nil {
		return g, errors.Errorf("could not find a license named %s", opts.Name)
	}
	g.File(genny.NewFile("LICENSE.plush", bytes.NewReader(body)))

	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))
	return g, nil
}
