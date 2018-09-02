package new

import (
	"errors"
)

type Options struct {
	Name string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("you must provide a Name")
	}
	return nil
}
