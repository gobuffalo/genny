package gogen

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/stretchr/testify/require"
)

func Test_AddGlobal(t *testing.T) {

	tt := []struct {
		Name     string
		Source   string
		Expected string
	}{
		{"single", globalBeforeSingle, globalAfterSingle},
		{"block", globalBeforeBlock, globalAfterBlock},
		{"empty1", globalBeforeEmpty1, globalAfterEmpty1},
		{"empty2", globalBeforeEmpty2, globalAfterEmpty2},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(st *testing.T) {
			r := require.New(t)
			path := filepath.Join("actions", "app.go")
			f := genny.NewFile(path, strings.NewReader(tc.Source))

			f, err := AddGlobal(f, "amount int", "top string")
			r.NoError(err)

			b, err := ioutil.ReadAll(f)
			r.NoError(err)

			r.Equal(tc.Expected, string(b))
		})
	}
}

const globalBeforeSingle = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var count int
var chkme int
func main() {
	var local int
	return
}
`

const globalAfterSingle = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var count int
var chkme int
var amount int
var top string
func main() {
	var local int
	return
}
`

const globalBeforeBlock = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var isOK bool

var (
	count int
	chkme int

)
func main() {
	var local int
	return
}
`

const globalAfterBlock = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var isOK bool

var (
	count int
	chkme int

	amount int
	top string
)
func main() {
	var local int
	return
}
`

const globalBeforeEmpty1 = `package actions

import "github.com/gobuffalo/buffalo"
import "github.com/gobuffalo/envy"

func main() {
	var local int
	return
}
`

const globalAfterEmpty1 = `package actions

import "github.com/gobuffalo/buffalo"
import "github.com/gobuffalo/envy"
var amount int
var top string

func main() {
	var local int
	return
}
`

const globalBeforeEmpty2 = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)
func main() {
	var local int
	return
}
`

const globalAfterEmpty2 = `package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)
var amount int
var top string
func main() {
	var local int
	return
}
`
