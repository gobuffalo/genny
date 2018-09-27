package gotools

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_ReplaceBlockContent(t *testing.T) {
	r := require.New(t)

	path := filepath.Join("actions", "app.go")
	f := genny.NewFile(path, strings.NewReader(modelBeforeBodyReplace))

	f, err := ReplaceBlockContent(f, "type X struct {", "Name string")
	r.NoError(err)

	b, err := ioutil.ReadAll(f)
	r.NoError(err)

	r.Equal(path, f.Name())
	r.Equal(modelAfterBodyReplace, string(b))
}

const modelBeforeBodyReplace = `
package models

type X struct {
	ID int
	Something string
}`

const modelAfterBodyReplace = `
package models

type X struct {
Name string
}`
