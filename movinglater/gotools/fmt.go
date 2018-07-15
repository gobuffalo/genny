package gotools

import (
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func GoFmt(root string) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		files, err := GoFiles(root)
		if err != nil {
			return errors.WithStack(err)
		}
		cmd, err := goFmtCmd(files...)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.Exec(cmd)
	})

	return g, nil
}

func goFmtCmd(files ...string) (*exec.Cmd, error) {
	if len(files) == 0 {
		files = []string{"."}
	}
	c := "gofmt"
	_, err := exec.LookPath("goimports")
	if err == nil {
		c = "goimports"
	}
	_, err = exec.LookPath("gofmt")
	if err != nil {
		return nil, errors.New("could not find gofmt or goimports")
	}
	args := []string{"-w"}
	args = append(args, files...)
	return exec.Command(c, args...), nil
}
