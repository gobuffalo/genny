package licenser

import (
	"github.com/gobuffalo/genny"
	gl "github.com/gobuffalo/licenser/genny/licenser"
	"github.com/markbates/oncer"
)

var Available = gl.Available

func New(opts *Options) (*genny.Generator, error) {
	oncer.Deprecate(0, "github.com/gobuffalo/genny/movinglater/licenser.New", "use github.com/gobuffalo/licenser/genny/licenser.New instead")
	return gl.New(opts)
}

type Options = gl.Options
