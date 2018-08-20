package gotools

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// AddImport adds n number of import statements into the path provided
func AddImport(gf genny.File, imports ...string) (genny.File, error) {
	name := gf.Name()
	gf, err := beforeParse(gf)
	if err != nil {
		return gf, errors.WithStack(err)
	}

	src, err := ioutil.ReadAll(gf)
	if err != nil {
		return gf, errors.WithStack(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gf.Name(), string(src), 0)
	if err != nil {
		return gf, errors.WithStack(err)
	}

	srcContent := string(src)
	fileLines := strings.Split(srcContent, "\n")

	end := findLastImport(f, fset, fileLines)

	x := make([]string, len(imports), len(imports)+2)
	for _, i := range imports {
		x = append(x, fmt.Sprintf("\t\"%s\"", i))

	}
	if end < 0 {
		x = append([]string{"import ("}, x...)
		x = append(x, ")")
	}

	fileLines = append(fileLines[:end], append(x, fileLines[end:]...)...)

	fileContent := strings.Join(fileLines, "\n")
	return genny.NewFile(name, strings.NewReader(fileContent)), err
}

func findLastImport(f *ast.File, fset *token.FileSet, fileLines []string) int {
	var end = -1

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			end = fset.Position(x.End()).Line
			return true
		}
		return true
	})

	return end
}
