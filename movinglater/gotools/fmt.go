package gotools

import (
	"os/exec"

	"github.com/pkg/errors"
)

func GoFmt(files ...string) (*exec.Cmd, error) {
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
