package gomods

import (
	"os/exec"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

const ENV = "GO111MODULE"

var ErrModsOff = errors.New("go mods are turned off")
var modsOn = (strings.TrimSpace(envy.Get(ENV, "off")) == "on")

func On() bool {
	return modsOn
}

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
	if !modsOn {
		return nil, ErrModsOff
	}
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		return r.Chdir(path, func() error {
			cmd := exec.Command(genny.GoBin(), "mod", "init", name)
			return r.Exec(cmd)
		})
	})
	return g, nil
}
