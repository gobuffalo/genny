package gotools

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools/goimports"
	"github.com/pkg/errors"
)

func GoFmt(root string) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		i, err := goimports.New(root)
		if err != nil {
			return errors.WithStack(err)
		}
		return i.Run()
	})

	return g, nil
}
