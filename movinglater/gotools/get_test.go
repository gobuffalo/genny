package gotools

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_WithGoGet(t *testing.T) {
	r := require.New(t)

	g := genny.Background()
	g = WithGoGet(g, "github.com/gobuffalo/envy", "-v")

	ca, ok := g.(genny.Commandable)
	r.True(ok)
	cmd := ca.Cmd()
	r.Equal("go get -v github.com/gobuffalo/envy", strings.Join(cmd.Args, " "))
}
