package gogen

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/gobuffalo/genny/v2"
)

// AddGlobal adds global variables.
func AddGlobal(gf genny.File, globals ...string) (genny.File, error) {
	pf, err := ParseFile(gf)
	if err != nil {
		return gf, err
	}

	gf = pf.File

	end, isGroup := findLastGlobal(pf.Ast, pf.FileSet, pf.Lines)
	if end < 0 {
		return gf, fmt.Errorf("unable to find the position to add the variables")
	}

	x := []string{}
	if isGroup {
		for _, gv := range globals {
			x = append(x, "\t"+gv)
		}
	} else {
		for _, gv := range globals {
			x = append(x, "var "+gv)
		}
	}

	pf.Lines = append(pf.Lines[:end], append(x, pf.Lines[end:]...)...)
	fileContent := strings.Join(pf.Lines, "\n")

	return genny.NewFile(gf.Name(), strings.NewReader(fileContent)), nil
}

func findLastGlobal(f *ast.File, fset *token.FileSet, fileLines []string) (int, bool) {
	var end = -1
	var isGroup = false

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok == token.VAR && x.Rparen > 0 {
				end = fset.Position(x.Rparen).Line - 1
				isGroup = true
				return false // if var block was found, break the loop
			} else if x.Tok == token.IMPORT {
				end = fset.Position(x.Rparen).Line
				if end == 0 {
					end = fset.Position(x.End()).Line
				} // this is the last resort. just after the import block
				return true
			}

		case *ast.ValueSpec:
			end = fset.Position(x.End()).Line
			return true // continue the block and find the last one

		case *ast.FuncDecl:
			if x.Name.String() == "main" {
				return false // break the loop if the main function started
			}
		}
		return true
	})

	return end, isGroup
}
