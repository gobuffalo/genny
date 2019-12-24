package gomods

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {

	table := []struct {
		name string
		root string
		mod  string
		err  bool
	}{
		{name: "coke", root: "", mod: "coke"},
		{name: "github.com/markbates/coke", root: "", mod: "github.com/markbates/coke"},
		{name: "", root: "", mod: "github.com/gobuffalo/genny/gogen/gomods"},
		{name: "", root: "../", mod: "github.com/gobuffalo/genny/gogen"},
	}

	for _, tt := range table {
		t.Run(fmt.Sprintf("%s/%s/%s", tt.name, tt.root, tt.mod), func(st *testing.T) {
			r := require.New(st)

			run := gentest.NewRunner()

			gg, err := New(tt.name, tt.root)
			r.NoError(err)
			run.WithGroup(gg)

			err = run.Run()
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			res := run.Results()
			r.Len(res.Files, 0)
			r.Len(res.Commands, 2)

			c := res.Commands[0]
			args := []string{genny.GoBin(), "mod", "init"}
			if len(tt.mod) > 0 {
				args = append(args, tt.mod)
			}
			r.Equal(args, c.Args)

			c = res.Commands[1]
			r.Equal([]string{genny.GoBin(), "mod", "tidy"}, c.Args)

		})
	}

	envy.Temp(func() {
		// envy.MustSet(ENV, "on")
		//
		// type ts struct {
		// 	name string
		// 	path string
		// 	mod  string
		// }
		// table := []ts{
		// 	{"", "", "github.com/gobuffalo/genny/gogen/gomods"},
		// 	{"coke", "", "coke"},
		// 	{"github.com\\gobuffalo\\coke", "", "github.com/gobuffalo/coke"},
		// 	{"", "github.com\\gobuffalo\\coke", "github.com/gobuffalo/coke"},
		// }
		//
		// c := build.Default
		//
		// for _, src := range c.SrcDirs() {
		// 	table = append(table, ts{
		// 		name: "coke",
		// 		path: src,
		// 		mod:  "coke",
		// 	})
		// 	table = append(table, ts{
		// 		name: "coke",
		// 		path: strings.ToLower(src),
		// 		mod:  "coke",
		// 	})
		// 	table = append(table, ts{
		// 		name: "",
		// 		path: filepath.Join(src, "github.com", "gobuffalo", "coke"),
		// 		mod:  "github.com/gobuffalo/coke",
		// 	})
		// }
		//
		// for _, tt := range table {
		// 	t.Run(tt.name+tt.path+tt.mod, func(st *testing.T) {
		// 		r := require.New(st)
		//
		// 		run := gentest.NewRunner()
		//
		// 		gg, err := New(tt.name, tt.path)
		// 		r.NoError(err)
		// 		run.WithGroup(gg)
		//
		// 		r.NoError(run.Run())
		// 		res := run.Results()
		// 		r.Len(res.Files, 0)
		// 		r.Len(res.Commands, 2)
		//
		// 		c := res.Commands[0]
		// 		args := []string{genny.GoBin(), "mod", "init"}
		// 		if len(tt.mod) > 0 {
		// 			args = append(args, tt.mod)
		// 		}
		// 		r.Equal(args, c.Args)
		//
		// 		c = res.Commands[1]
		// 		r.Equal([]string{genny.GoBin(), "mod", "tidy"}, c.Args)
		// 	})
		// }
		//
	})
}
