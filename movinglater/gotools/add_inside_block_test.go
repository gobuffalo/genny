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

func Test_AddInsideBlock(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	dir := os.TempDir()
	err := run.Chdir(dir, func() error {
		path := filepath.Join("actions", "app.go")
		tf, err := os.Create(path)
		r.NoError(err)
		tf.Write([]byte(appBefore))
		tf.Close()

		f := genny.NewFile(path, strings.NewReader(appBefore))
		f, err = AddInsideBlock(f, "if app == nil {", "app.Use(Foo)")
		r.NoError(err)

		b, err := ioutil.ReadAll(f)
		r.NoError(err)

		r.Equal(appAfter, string(b))
		return nil
	})
	r.NoError(err)
}

func Test_AddInsideBlock_NoFound(t *testing.T) {
	r := require.New(t)

	run := genny.DryRunner(context.Background())
	dir := os.TempDir()
	err := run.Chdir(dir, func() error {
		path := filepath.Join("actions", "app.go")
		tf, err := os.Create(path)
		r.NoError(err)
		tf.Write([]byte(appBefore))
		tf.Close()

		f := genny.NewFile(path, strings.NewReader(appBefore))
		f, err = AddInsideBlock(f, "idontexist", "app.Use(Foo)")
		r.Error(err)

		return nil
	})
	r.NoError(err)
}

const appBefore = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/mw-paramlogger"

	"github.com/gobuffalo/mw-csrf"
	"github.com/gobuffalo/mw-poptx"
	"github.com/markbates/coke/models"
)

var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
		app.Use(paramlogger.ParameterLogger)
		app.Use(csrf.New)
		app.Use(poptx.PopTransaction(models.DB))
		app.GET("/", HomeHandler)
		app.ServeFiles("/", assetsBox)
	}

	return app
}`

const appAfter = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/mw-paramlogger"

	"github.com/gobuffalo/mw-csrf"
	"github.com/gobuffalo/mw-poptx"
	"github.com/markbates/coke/models"
)

var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{})
		app.Use(paramlogger.ParameterLogger)
		app.Use(csrf.New)
		app.Use(poptx.PopTransaction(models.DB))
		app.GET("/", HomeHandler)
		app.Use(Foo)
		app.ServeFiles("/", assetsBox)
	}

	return app
}`
