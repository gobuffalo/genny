package new

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	r.Error(opts.Validate())

	opts.Name = "foo"
	r.NoError(opts.Validate())

	r.Equal("github.com/gobuffalo/genny/foo/templates", opts.BoxName)
}
