package genny

import (
	"context"
	"io/ioutil"
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

	run := DryRunner(context.Background())
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 1)

	c := res.Commands[0]
	r.Equal("foo bar", strings.Join(c.Args, " "))

	f := res.Files[0]
	r.Equal("foo.txt", f.Name())
	r.Equal("Hello mark", f.String())
}

func Test_Runner_FindFile(t *testing.T) {
	r := require.New(t)

	g := New()
	g.File(NewFile("foo.txt", strings.NewReader("Hello mark")))
	g.File(NewFile("foo.txt", strings.NewReader("Hello world")))

	run := DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 1)

	f, err := run.FindFile("foo.txt")
	r.NoError(err)
	r.Equal(res.Files[0], f)
}

func Test_Runner_FindFile_FromDisk(t *testing.T) {
	r := require.New(t)

	run := DryRunner(context.Background())

	exp, err := ioutil.ReadFile("./fixtures/foo.txt")
	r.NoError(err)

	f, err := run.FindFile("fixtures/foo.txt")
	r.NoError(err)
	act, err := ioutil.ReadAll(f)
	r.NoError(err)

	r.Equal(string(exp), string(act))
}

func Test_Runner_FindFile_DoesntExist(t *testing.T) {
	r := require.New(t)

	run := DryRunner(context.Background())

	_, err := run.FindFile("idontexist")
	r.Error(err)
}
