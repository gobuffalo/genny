package new

import (
	"errors"
	"path"
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
		opts.BoxName = path.Join(opts.Prefix, opts.Name, "templates")
	}
	return nil
}
