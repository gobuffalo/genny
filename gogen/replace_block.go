package gogen

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny/v2"
)

// ReplaceBlock replaces found block with given new open/close statements
/*
	open - start statement of the block e.g. `if isTrue {`
	close - close statement of the block e.g. `}`
	newOpen - new start statement of the block e.g. `defer func(isTrue) {`
	newClose - new close statement of the block e.g. `}()`
*/
func ReplaceBlock(gf genny.File, open, close, newOpen, newClose string) (genny.File, error) {
	pf, err := ParseFile(gf)
	if err != nil {
		return gf, err
	}
	gf = pf.File

	start, end := findBlockCoordinates(open, pf)
	if end < 0 {
		return gf, fmt.Errorf("could not find desired block in %s", gf.Name())
	}

	lines := pf.Lines[:start-1]
	lines = append(lines, strings.Replace(pf.Lines[start-1], open, newOpen, 1)) // open
	lines = append(lines, pf.Lines[start:end]...)                               // body
	lines = append(lines, strings.Replace(pf.Lines[end], close, newClose, 1))   // close
	lines = append(lines, pf.Lines[end+1:]...)
	fileContent := strings.Join(lines, "\n")

	return genny.NewFile(gf.Name(), strings.NewReader(fileContent)), nil
}
