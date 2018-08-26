package genny

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_Generator_File(t *testing.T) {
	r := require.New(t)

	g := New()
	g.File(NewFile("foo.txt", strings.NewReader("hello")))

	run := DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("foo.txt", f.Name())
	r.Equal("hello", f.String())
}

func Test_Generator_Box(t *testing.T) {
	r := require.New(t)

	g := New()
	r.NoError(g.Box(packr.NewBox("./fixtures")))

	run := DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("bar/baz.txt", f.Name())
	r.Equal("baz!\n", f.String())

	f = res.Files[1]
	r.Equal("foo.txt", f.Name())
	r.Equal("foo!\n", f.String())
}

func Test_Command(t *testing.T) {
	r := require.New(t)

	g := New()
	g.Command(exec.Command("echo", "hello"))

	run := DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 0)

	c := res.Commands[0]
	r.Equal("echo hello", strings.Join(c.Args, " "))
}
