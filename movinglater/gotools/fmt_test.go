package gotools

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_GoFmt(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())

	g, err := GoFmt("")
	r.NoError(err)
	run.With(g)

	r.NoError(run.Run())

}
