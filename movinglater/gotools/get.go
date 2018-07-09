package gotools

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

// WithGoGet returns a `go get` command
func WithGoGet(g genny.Generator, pkg string, args ...string) genny.Generator {
	args = append([]string{"get"}, args...)
	args = append(args, pkg)
	cmd := exec.CommandContext(g.Context(), "go", args...)
	g = genny.WithCmd(g, cmd)
	return g
}
