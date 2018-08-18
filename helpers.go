package genny

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

func exts(f File) []string {
	var exts []string

	name := f.Name()
	ext := filepath.Ext(name)

	for ext != "" {
		exts = append([]string{ext}, exts...)
		name = strings.TrimSuffix(name, ext)
		ext = filepath.Ext(name)
	}
	return exts
}

// HasExt checks if a file has a particular extension
func HasExt(f File, ext string) bool {
	if ext == "*" {
		return true
	}
	for _, x := range exts(f) {
		if x == ext {
			return true
		}
	}
	return false
}

// StripExt from a File and return a new one
func StripExt(f File, ext string) File {
	name := f.Name()
	name = strings.Replace(name, ext, "", -1)
	return NewFile(name, f)
}

// Chdir will change to the specified directory
// and revert back to the current directory when
// the runner function has returned.
// If the directory does not exist, it will be
// created for you.
func Chdir(path string, fn func() error) error {
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	os.MkdirAll(path, 0755)
	if err := os.Chdir(path); err != nil {
		return errors.WithStack(err)
	}
	if err := fn(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GoBin() string {
	return envy.Get("GO_BIN", "go")
}
