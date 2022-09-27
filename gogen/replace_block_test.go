package gogen

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/stretchr/testify/require"
)

func Test_ReplaceBlock(t *testing.T) {
	r := require.New(t)

	path := "actions/app.go"
	f := genny.NewFile(path, strings.NewReader(modelBeforeBlockReplace))

	f, err := ReplaceBlock(f, "if app == nil {", "}", "appOnce.Do(func() {", "})")

	r.NoError(err)

	b, err := ioutil.ReadAll(f)
	r.NoError(err)

	r.Equal(path, f.Name())
	r.Equal(modelAfterBlockReplace, string(b))
}

const modelBeforeBlockReplace = `
package actions

import (
	"sync"

	"github.com/gobuffalo/buffalo"
)

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
		app.GET("/", HomeHandler)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}`

const modelAfterBlockReplace = `
package actions

import (
	"sync"

	"github.com/gobuffalo/buffalo"
)

func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{})
		app.GET("/", HomeHandler)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	})

	return app
}`
