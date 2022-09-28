package gogen

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/gobuffalo/genny/v2"
)

// AddInsideBlock will add anything inside of the app declaration block inside of file
func AddInsideBlock(gf genny.File, search string, expressions ...string) (genny.File, error) {
	pf, err := ParseFile(gf)
	if err != nil {
		return gf, err
	}
	gf = pf.File

	_, end := findBlockCoordinates(search, pf)
	if end < 0 {
		return gf, fmt.Errorf("could not find desired block in %s", gf.Name())
	}

	if len(pf.Lines) == end {
		end = end - 1
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

func findBlockCoordinates(search string, pf ParsedFile) (int, int) {
	var end = -1
	var start = -1

	ast.Inspect(pf.Ast, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType, *ast.BlockStmt:
			line := pf.FileSet.Position(x.Pos()).Line
			structDeclaration := fmt.Sprintf("%s\n", pf.Lines[line-1])

			if strings.Contains(structDeclaration, search) {
				start = line
				end = pf.FileSet.Position(x.End()).Line - 1
				return false // should return false to guarantee the result
			}
		}
		return true
	})

	return start, end
}
