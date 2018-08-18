package dep

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func Init(path string, verbose bool) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		return genny.Chdir(path, func() error {
			if _, err := exec.LookPath("dep"); err != nil {
				return err
			}
			cmd := exec.Command("dep", "init")
			if verbose {
				cmd.Args = append(cmd.Args, "-v")
			}
			return r.Exec(cmd)
		})
	})
	return g, nil
}
