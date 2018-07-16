package licenser

import (
	"bytes"
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/sirupsen/logrus"
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

	bb := &bytes.Buffer{}
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.Out = bb

	run := genny.DryRunner(context.Background())
	run.Logger = l
	run.With(g)

	r.NoError(run.Run())

	out := bb.String()
	r.Contains(out, "LICENSE")
	r.Contains(out, "Apache License")
	r.NotContains(out, "The MIT License")
}

func Test_Available(t *testing.T) {
	r := require.New(t)
	r.Len(Available, 16)
}
