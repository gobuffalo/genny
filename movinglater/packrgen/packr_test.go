package packrgen

import (
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_WithBox(t *testing.T) {
	r := require.New(t)
	box := packr.NewBox("./fixtures")

	var names []string
	g := genny.Background()
	g, err := WithBox(g, box, func(gg genny.Generator, f genny.File) genny.Generator {
		names = append(names, f.Name())
		return genny.WithFile(gg, f)
	})
	r.NoError(err)
	r.NoError(g.Run())

	r.Len(names, 2)
}
