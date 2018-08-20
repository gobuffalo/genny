package gotools

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_AddImport(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	dir := os.TempDir()
	err := run.Chdir(dir, func() error {
		path := filepath.Join("actions", "app.go")
		tf, err := os.Create(path)
		r.NoError(err)
		tf.Write([]byte(importBefore))
		tf.Close()

		f := genny.NewFile(path, strings.NewReader(importBefore))
		f, err = AddImport(f, "foo/bar", "foo/baz")
		r.NoError(err)

		b, err := ioutil.ReadAll(f)
		r.NoError(err)

		r.Equal(importAfter, string(b))
		return nil
	})
	r.NoError(err)
}

const importBefore = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)
`

const importAfter = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"


	"foo/bar"
	"foo/baz"
)
`
