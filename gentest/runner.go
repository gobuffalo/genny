package gentest

import (
	"context"

	"github.com/gobuffalo/genny"
)

func NewRunner() *genny.Runner {
	r := genny.DryRunner(context.Background())
	r.Logger = NewLogger()
	return r
}
