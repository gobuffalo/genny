package gomods

import (
	"errors"
	"os"

	"github.com/gobuffalo/genny/internal/takeon/github.com/markbates/safe"
	"github.com/gobuffalo/here"
)

const ENV = "GO111MODULE"

var ErrModsOff = errors.New("go mods are turned off")

func Force(b bool) {
	if b {
		os.Setenv(ENV, "on")
		return
	}
	os.Setenv(ENV, "off")
}

func On() bool {
	oe := os.Getenv(ENV)
	if oe == "off" {
		return false
	}

	info, _ := here.Current()
	return !info.Module.IsZero()
}

func Disable(fn func() error) error {
	oe := os.Getenv(ENV)
	if oe == "" {
		oe = "off"
	}

	os.Setenv(ENV, "off")

	err := safe.RunE(fn)
	os.Setenv(ENV, oe)
	return err
}
