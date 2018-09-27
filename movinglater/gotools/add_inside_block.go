package gotools

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// AddInsideBlock will add anything inside of the app declaration block inside of file
func AddInsideBlock(gf genny.File, search string, expressions ...string) (genny.File, error) {
	pf, err := ParseFile(gf)
	if err != nil {
		return gf, errors.WithStack(err)
	}
	gf = pf.File

	end := findClosingRouteBlockEnd(search, pf.Ast, pf.FileSet, pf.Lines)
	if end < 0 {
		return gf, errors.Errorf("could not find desired block in %s", gf.Name())
	}

	el := pf.Lines[end:]
	sl := []string{}
	sf := []string{}
	for _, l := range pf.Lines[:end] {
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
	pf.Lines = append(sl, append(expressions, el...)...)

	fileContent := strings.Join(pf.Lines, "\n")
	return genny.NewFile(gf.Name(), strings.NewReader(fileContent)), nil
}

func findClosingRouteBlockEnd(search string, f *ast.File, fset *token.FileSet, fileLines []string) int {
	var end = -1

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			line := fset.Position(x.Pos()).Line
			structDeclaration := fmt.Sprintf("%s\n", fileLines[line-1])

			if strings.Contains(structDeclaration, search) {
				end = fset.Position(x.End()).Line - 1
				return false
			}

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

func findBlockCoordinates(search string, f *ast.File, fset *token.FileSet, fileLines []string) (int, int) {
	var end = -1
	var start = -1

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			line := fset.Position(x.Pos()).Line
			structDeclaration := fmt.Sprintf("%s\n", fileLines[line-1])

			if strings.Contains(structDeclaration, search) {
				start = line
				end = fset.Position(x.End()).Line
				return false
			}

		case *ast.BlockStmt:
			start = fset.Position(x.Lbrace).Line
			blockDeclaration := fmt.Sprintf("%s\n", fileLines[start-1])

			if strings.Contains(blockDeclaration, search) {
				end = fset.Position(x.Rbrace).Line - 1
			}

		}
		return true
	})

	return start, end
}
