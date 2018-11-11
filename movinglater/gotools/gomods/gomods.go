package gomods

import (
	"github.com/gobuffalo/envy"
	"github.com/markbates/safe"
	"github.com/pkg/errors"
)

const ENV = "GO111MODULE"

var ErrModsOff = errors.New("go mods are turned off")

func Force(b bool) {
	if b {
		envy.Set(ENV, "on")
		return
	}
	envy.Set(ENV, "off")
}

func On() bool {
	return envy.Mods()
}

func Disable(fn func() error) error {
	var err error
	envy.Temp(func() {
		envy.MustSet(ENV, "off")
		err = safe.RunE(fn)
	})
	return err
}
