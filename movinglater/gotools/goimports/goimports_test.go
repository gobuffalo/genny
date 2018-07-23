package goimports

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_processFile(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}
	i := NewFromFiles(File{
		Name: "foo.go",
		In:   strings.NewReader(tmpl),
		Out:  bb,
	})
	r.NoError(i.Run())
	r.Equal(actual, bb.String())
}

const tmpl = `package foo

func init() {
	packr.NewBox("")
}
`

const actual = `package foo

import "github.com/gobuffalo/packr"

func init() {
	packr.NewBox("")
}
`
