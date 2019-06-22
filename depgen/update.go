package depgen

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func Update(verbose bool) (*genny.Generator, error) {
	g := genny.New()
	var args []string
	if verbose {
		args = append(args, "-v")
	}

	id, err := InstallDep(args...)
	if err != nil {
		return g, err
	}
	g.Merge(id)

	cmd := exec.Command("dep", "ensure")
	args = append(args, "-update")
	if verbose {
		cmd.Args = append(cmd.Args, args...)
	}
	return g, nil
}
