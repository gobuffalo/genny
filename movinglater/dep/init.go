package dep

import (
	"os"
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func Init(path string, verbose bool) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		defer os.Chdir(pwd)
		if err := os.Chdir(path); err != nil {
			return errors.WithStack(err)
		}

		if _, err := exec.LookPath("dep"); err != nil {
			return err
		}
		cmd := exec.Command("dep", "init")
		if verbose {
			cmd.Args = append(cmd.Args, "-v")
		}
		return r.Exec(cmd)
	})
	return g, nil
}
