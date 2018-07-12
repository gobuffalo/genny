package genny

import (
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

	run, bb := testRunner()
	run.With(g)
	r.NoError(run.Run())

	out := bb.String()
	r.Contains(out, "foo.txt")
	r.Contains(out, "hello")
}

func Test_Generator_Box(t *testing.T) {
	r := require.New(t)

	g := New()
	r.NoError(g.Box(packr.NewBox("./fixtures")))

	run, bb := testRunner()
	run.With(g)
	r.NoError(run.Run())

	out := bb.String()
	r.Contains(out, "bar/baz.txt")
	r.Contains(out, "baz!")
	r.Contains(out, "foo.txt")
	r.Contains(out, "foo!")
}

func Test_Command(t *testing.T) {
	r := require.New(t)

	g := New()
	g.Command(exec.Command("echo", "hello"))

	run, bb := testRunner()
	run.With(g)
	r.NoError(run.Run())

	out := bb.String()
	r.Contains(out, "echo hello")
}
