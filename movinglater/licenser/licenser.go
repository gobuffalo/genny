package licenser

import (
	gl "github.com/gobuffalo/licenser/genny/licenser"
	"github.com/markbates/oncer"
)

func init() {
	oncer.Deprecate(5, "github.com/gobuffalo/genny/movinglater/licenser", "use github.com/gobuffalo/licenser/genny/licenser instead")
}

var Available = gl.Available

var New = gl.New

type Options = gl.Options
