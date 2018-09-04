package genny

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
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
		FileFn: func(f File) (File, error) {
			defer func() {
				if s, ok := f.(io.Seeker); ok {
					s.Seek(0, 0)
				}
			}()
			if _, err := io.Copy(os.Stdout, f); err != nil {
				return f, errors.WithStack(err)
			}
			return f, nil
		},
	}
	r.Disk = newDisk(r)
	return r
}
