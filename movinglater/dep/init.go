package dep

import (
	"os"
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func Init(path string, verbose bool) (*genny.Generator, error) {
	g := genny.New()
	pwd, err := os.Getwd()
	if err != nil {
		return g, errors.WithStack(err)
	}
	defer os.Chdir(pwd)
	if err := os.Chdir(path); err != nil {
		return g, errors.WithStack(err)
	}

	if _, err := exec.LookPath("dep"); err != nil {
		return g, nil
	}
	cmd := exec.Command("dep", "init")
	if verbose {
		cmd.Args = append(cmd.Args, "-v")
	}
	g.Command(cmd)
	return g, nil
}
