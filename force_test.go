package genny

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Force_Exists(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "test")
	r.NoError(err)

	g, err := Force(dir, false)
	r.NoError(err)

	run := DryRunner(context.Background())
	run.With(g)
	r.Error(run.Run())

	g, err = Force(dir, true)
	r.NoError(err)

	run = DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())
}

func Test_Force_Doesnt_Exists(t *testing.T) {
	r := require.New(t)

	dir := "i don't exist"
	g, err := Force(dir, false)
	r.NoError(err)

	run := DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	g, err = Force(dir, true)
	r.NoError(err)

	run = DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())
}
