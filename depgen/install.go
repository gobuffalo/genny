package depgen

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func InstallDep(args ...string) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		if _, err := r.LookPath("dep"); err == nil {
			return nil
		}

		args = append([]string{"get"}, args...)
		args = append(args, "github.com/golang/dep/cmd/dep")
		c := exec.Command(genny.GoBin(), args...)
		return r.Exec(c)
	})
	return g, nil
}
