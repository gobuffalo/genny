package genny

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

// DryRun will wrap the Generator with `DryExec` to prevent
// cmds from executing, as well as a FilesHandler that will
// echo file contents to the screen instead of disk.
func DryRun(g Generator) Generator {
	g = DryExec(g)
	g = WithFilesHandler(g, func(f File) error {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.WithStack(err)
		}
		g.Logger().Printf(string(b))
		return nil
	})
	return g
}
