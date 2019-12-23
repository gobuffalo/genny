package new

import (
	"errors"
	"path"

	"github.com/gobuffalo/here"
)

type Options struct {
	Prefix  string
	Name    string
	BoxName string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("you must provide a Name")
	}

	if len(opts.BoxName) == 0 {
		info, err := here.Current()
		if err != nil {
			return err
		}
		opts.BoxName = path.Join(info.ImportPath, opts.Prefix, opts.Name, "templates")
	}
	return nil
}
