package new

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{Name: "Foo"})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 4)

	f := res.Files[0]
	r.Equal("foo/foo.go", f.Name())
	body := f.String()
	r.Contains(body, "package foo")
	r.Contains(body, "../foo/templates")

	f = res.Files[1]
	r.Equal("foo/foo_test.go", f.Name())
	body = f.String()
	r.Contains(body, "package foo")

	f = res.Files[2]
	r.Equal("foo/options.go", f.Name())
	body = f.String()
	r.Contains(body, "package foo")

	f = res.Files[3]
	r.Equal("foo/templates/example.txt", f.Name())
}
