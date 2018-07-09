package genny

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// WetFilesHandler returns a generator that will write any
// previously defined files to disk. It was also create any
// directories as needed
func WetFilesHandler(g Generator) Generator {
	return WithFilesHandler(g, func(file File) error {
		g.Logger().Infof("[genny.WetFilesHandler] %s\n", file.Name())
		dir := filepath.Dir(file.Name())
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.WithStack(err)
		}
		f, err := os.Create(file.Name())
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}
