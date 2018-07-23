package genny

import (
	"os"

	"github.com/pkg/errors"
)

// Force will return an error if the path exists if `force` is false. If `force` is true it will delete the path.
func Force(path string, force bool) (*Generator, error) {
	g := New()
	g.RunFn(func(r *Runner) error {
		_, err := os.Stat(path)
		if err != nil {
			// path doesn't exist. move on.
			return nil
		}
		if !force {
			return errors.Errorf("path %s already exists", path)
		}
		if err := os.RemoveAll(path); err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	return g, nil
}
