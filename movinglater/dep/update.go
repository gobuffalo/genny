package dep

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func Update(verbose bool) (*genny.Generator, error) {
	g := genny.New()
	if _, err := exec.LookPath("dep"); err != nil {
		return g, nil
	}
	cmd := exec.Command("dep", "ensure", "-update")
	if verbose {
		cmd.Args = append(cmd.Args, "-v")
	}
	g.Command(cmd)
	return g, nil
}
