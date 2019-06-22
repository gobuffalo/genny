package gogen

import (
	"path/filepath"

	"errors"

	"github.com/gobuffalo/genny"
)

func PackageName(f genny.File) (string, error) {
	pkg := filepath.Base(filepath.Dir(f.Name()))
	pf, err := ParseFile(f)
	if err == nil {
		pkg = pf.Ast.Name.String()
	}
	if len(pkg) == 0 || pkg == "." {
		return "", errors.New("could not determine package")
	}
	return pkg, nil
}
