package command

import (
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

type Options struct {
	App    meta.App
	Prefix string
	Name   string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	if len(opts.Name) == 0 {
		return errors.New("you must provide a name")
	}
	return nil
}
