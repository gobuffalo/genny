package gomods

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New_No_Modules(t *testing.T) {
	om := modsOn
	defer func() { modsOn = om }()
	modsOn = false

	r := require.New(t)
	_, err := New("foo/bar", "")
	r.Error(err)
}

func Test_New_With_Modules(t *testing.T) {
	om := modsOn
	defer func() { modsOn = om }()
	modsOn = true
	r := require.New(t)
	run := genny.DryRunner(context.Background())

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)
	g, err := New("foo/bar", dir)
	r.NoError(err)
	g.With(run)

	r.NoError(run.Run())
	res := run.Results()
	r.Len(res.Files, 0)
	r.Len(res.Commands, 2)

	c := res.Commands[0]
	r.Equal([]string{genny.GoBin(), "mod", "init", "foo/bar"}, c.Args)

	c = res.Commands[1]
	r.Equal([]string{genny.GoBin(), "mod", "tidy"}, c.Args)
}
