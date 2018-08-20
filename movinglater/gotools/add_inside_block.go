package gotools

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// AddInsideBlock will add anything inside of the app declaration block inside of file
func AddInsideBlock(gf genny.File, search string, expressions ...string) (genny.File, error) {
	name := gf.Name()
	gf, err := beforeParse(gf)
	if err != nil {
		return gf, errors.WithStack(err)
	}

	src := gf.String()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gf.Name(), src, 0)
	if err != nil {
		return gf, err
	}

	fileLines := strings.Split(src, "\n")

	end := findClosingRouteBlockEnd(search, f, fset, fileLines)
	if end < 0 {
		return gf, errors.Errorf("could not find desired block in %s", name)
	}

	el := fileLines[end:]
	sl := []string{}
	sf := []string{}
	for _, l := range fileLines[:end] {
		// if there's a app.ServeFiles("/", foo) line it needs to be the last added to the router
		if strings.Contains(l, "ServeFiles(\"/\"") {
			sf = append(sf, l)
			continue
		}
		sl = append(sl, l)
	}

	for i := 0; i < len(expressions); i++ {
		expressions[i] = fmt.Sprintf("\t\t%s", expressions[i])
	}

	el = append(sf, el...)
	fileLines = append(sl, append(expressions, el...)...)

	fileContent := strings.Join(fileLines, "\n")
	return genny.NewFile(name, strings.NewReader(fileContent)), err
}

func findClosingRouteBlockEnd(search string, f *ast.File, fset *token.FileSet, fileLines []string) int {
	var end = -1

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BlockStmt:
			start := fset.Position(x.Lbrace).Line
			blockDeclaration := fmt.Sprintf("%s\n", fileLines[start-1])

			if strings.Contains(blockDeclaration, search) {
				end = fset.Position(x.Rbrace).Line - 1
			}

		}
		return true
	})

	return end
}

func beforeParse(gf genny.File) (genny.File, error) {
	src, err := ioutil.ReadAll(gf)
	if err != nil {
		return gf, errors.WithStack(err)
	}

	dir := os.TempDir()
	path := filepath.Join(dir, gf.Name())
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return gf, errors.WithStack(err)
	}

	tf, err := os.Create(path)
	if err != nil {
		return gf, errors.WithStack(err)
	}
	if _, err := tf.Write(src); err != nil {
		return gf, errors.WithStack(err)
	}
	if err := tf.Close(); err != nil {
		return gf, errors.WithStack(err)
	}
	return genny.NewFile(path, bytes.NewReader(src)), nil
}
