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
	bb := testLogger(run)

	g := New()
	g.Command(exec.Command("echo", "hello"))
	g.File(NewFile("foo.txt", strings.NewReader("foo!")))
	run.With(g)

	r.NoError(run.Run())
	r.Contains(bb.String(), "hello")

	b, err := ioutil.ReadFile(filepath.Join(run.Root, "foo.txt"))
	r.NoError(err)
	r.Equal("foo!", string(b))
}
