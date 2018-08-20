package genny

import (
	"context"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// DryRunner will NOT execute commands and write files
// it is NOT destructive
func DryRunner(ctx context.Context) *Runner {
	pwd, _ := os.Getwd()
	l := logrus.New()
	l.Out = os.Stdout
	l.SetLevel(logrus.DebugLevel)
	r := &Runner{
		Logger:  l,
		Context: ctx,
		Root:    pwd,
		moot:    &sync.RWMutex{},
		files:   map[string]File{},
	}
	return r
}
