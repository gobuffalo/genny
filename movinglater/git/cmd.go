package git

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

var ErrWorkingTreeClean = errors.New("working tree clean")

func Run(args ...string) genny.RunFn {
	return func(r *genny.Runner) error {
		bb := &bytes.Buffer{}
		cmd := exec.Command("git", args...)
		cmd.Stdout = bb
		cmd.Stderr = bb
		err := r.Exec(cmd)
		if err != nil {
			if strings.Contains(bb.String(), "working tree clean") {
				return ErrWorkingTreeClean
			}
			return errors.WithStack(err)
		}
		return nil
	}

}
