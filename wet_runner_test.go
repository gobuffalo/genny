package genny

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WetRunner(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	run := WetRunner(context.Background())
	run.Root = dir

	g := New()
	g.Command(exec.Command("echo", "hello"))
	g.File(NewFile("foo.txt", strings.NewReader("foo!")))
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 1)

	c := res.Commands[0]
	r.Equal("echo hello", strings.Join(c.Args, " "))

	f := res.Files[0]
	r.Equal("foo.txt", f.Name())
	r.Equal("foo!", f.String())

	b, err := ioutil.ReadFile(filepath.Join(run.Root, "foo.txt"))
	r.NoError(err)
	r.Equal("foo!", string(b))
}
