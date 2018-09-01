package licenser

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_Licenser(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		Year:   1999,
		Author: "Homer Simpson",
		Name:   "apache",
	}

	g, err := New(opts)
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("LICENSE", f.Name())
	body := f.String()
	r.Contains(body, "Apache License")
	r.NotContains(body, "The MIT License")
}

func Test_Available(t *testing.T) {
	r := require.New(t)
	r.Len(Available, 16)
}
