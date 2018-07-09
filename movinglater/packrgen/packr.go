package packrgen

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// FileFn for use with `WithBox`
type FileFn func(genny.Generator, genny.File) genny.Generator

// WithBox will walk through a packr.Box and call the *optional* FileFn allowing
// you to customize the returned generator for each file. The default, if this is `nil`
// will be `genny.WithFile` for each file in the box.
func WithBox(g genny.Generator, box packr.Box, fn FileFn) (genny.Generator, error) {
	if fn == nil {
		fn = func(g genny.Generator, f genny.File) genny.Generator { return genny.WithFile(g, f) }
	}
	err := box.Walk(func(path string, pf packr.File) error {
		fi, err := pf.FileInfo()
		if err != nil {
			return errors.WithStack(err)
		}
		if fi.IsDir() {
			return nil
		}
		f := genny.NewFile(path, pf)
		g = fn(g, f)
		return nil
	})
	return g, err
}
