package gotools

import (
	"os/exec"
	"strings"

	"github.com/gobuffalo/genny"
)

// WithGoGet returns a `go get` command
func WithGoGet(g genny.Generator, pkg string, args ...string) genny.Generator {
	args = append([]string{"get"}, args...)
	args = append(args, pkg)
	cmd := exec.CommandContext(g.Context(), "go", args...)
	key := strings.Join(cmd.Args, " ")
	for _, t := range genny.Tree(g) {
		if ca, ok := t.(genny.Commandable); ok {
			if ca.Cmd() != nil {
				cakey := strings.Join(ca.Cmd().Args, " ")
				if cakey == key {
					return g
				}
			}
		}
	}
	g = genny.WithCmd(g, cmd)
	return g
}
