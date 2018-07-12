package genny

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// WetRunner will execute commands and write files
// it is DESTRUCTIVE
func WetRunner(ctx context.Context) *Runner {
	r := DryRunner(ctx)
	l := logrus.New()
	l.Out = os.Stdout
	r.Logger = l

	r.ExecFn = wetExecFn
	r.FileFn = func(f File) error {
		return wetFileFn(r, f)
	}
	return r
}

func wetExecFn(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func wetFileFn(r *Runner, f File) error {
	name := f.Name()
	if !filepath.IsAbs(name) {
		name = filepath.Join(r.Root, name)
	}
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}
	ff, err := os.Create(name)
	if err != nil {
		return errors.WithStack(err)
	}
	defer ff.Close()
	if _, err := io.Copy(ff, f); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
