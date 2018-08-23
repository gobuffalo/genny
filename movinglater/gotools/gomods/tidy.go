package gomods

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func Tidy(path string, verbose bool) (*genny.Generator, error) {
	if !modsOn {
		return nil, ErrModsOff
	}
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
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
