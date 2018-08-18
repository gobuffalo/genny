package gomods

import (
	"os/exec"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

const ENV = "GO111MODULE"

func New(name string, path string) (*genny.Group, error) {
	g := &genny.Group{}

	init, err := Init(name, path)
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.Add(init)

	tidy, err := Tidy(path, false)
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.Add(tidy)
	return g, nil
}

func Init(name string, path string) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		if envy.Get(ENV, "") != "on" {
			return nil
		}
		return r.Chdir(path, func() error {
			cmd := exec.Command(genny.GoBin(), "mod", "init", name)
			return r.Exec(cmd)
		})
	})
	return g, nil
}

func Tidy(path string, verbose bool) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		if envy.Get(ENV, "") != "on" {
			return nil
		}
		return r.Chdir(path, func() error {
			cmd := exec.Command(genny.GoBin(), "mod", "tidy")
			if verbose {
				cmd.Args = append(cmd.Args, "-v")
			}
			return r.Exec(cmd)
		})
	})
	return g, nil
}
