package genny

import (
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Runner_Run(t *testing.T) {
	r := require.New(t)

	g := New()

	g.Command(exec.Command("foo", "bar"))
	g.File(NewFile("foo.txt", strings.NewReader("Hello mark")))

	run, bb := testRunner()
	run.FileFn = func(f File) error {
		io.Copy(bb, f)
		return nil
	}
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 1)

	out := bb.String()
	r.Contains(out, "foo bar")
	r.Contains(out, "foo.txt")
	r.Contains(out, "Hello mark")
}
