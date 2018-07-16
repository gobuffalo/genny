package licenser

import (
	"os/user"
	"strings"
	"time"

	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

type Options struct {
	Name   string
	Year   int
	Author string
}

func NormalizeOptions(opts *Options) error {
	opts.Name = defaults.String(opts.Name, "mit")
	opts.Name = strings.ToLower(opts.Name)
	if opts.Year == 0 {
		opts.Year = time.Now().Year()
	}

	if opts.Author == "" {
		u, err := user.Current()
		if err != nil {
			return errors.WithStack(err)
		}
		opts.Author = defaults.String(u.Name, u.Username)
	}

	return nil
}
