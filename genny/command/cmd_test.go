package command

import (
	"testing"

	"github.com/gobuffalo/genny/gentest"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:   "widget",
		Prefix: "genny",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	// r.Len(res.Commands, 0)
	// r.Len(res.Files, 3)

	box := packr.New("genny/cmd#Test_New", "./_fixtures")
	err = gentest.CompareBoxStripped(box, res)
	r.NoError(err)
	// f, err := res.Find("main.go")
	// r.NoError(err)
	// r.Contains(f.String(), `import "github.com/gobuffalo/genny/genny/cmd/widget/cmd"`)
}
