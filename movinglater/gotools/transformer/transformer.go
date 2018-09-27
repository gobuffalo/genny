package transformer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

//Transformer has helper classes to modify go source code
type Transformer struct {
	filePath string
	file     os.File
}

//NewTransformer creates a Transformer for a provided path of a file
func NewTransformer(path string) *Transformer {
	return &Transformer{
		filePath: path,
	}
}

//AddImports allows to add imports to a .go file
func (tr *Transformer) AddImports(imports ...string) error {
	content, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, tr.filePath, string(content), 0)
	if err != nil {
		return err
	}

	var end = 1
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			end = fset.Position(x.End()).Line
			return true
		}
		return true
	})

	lines := strings.Split(string(content), "\n")
	c := append(lines[:end], append(imports, lines[end:]...)...)
	fileContent := strings.Join(c, "\n")

	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//Append appends at the bottom of the source file
func (tr *Transformer) Append(src string) error {
	content, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	c := append(lines, strings.Split(src, "\n")...)
	fileContent := strings.Join(c, "\n")

	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//AppendAfter Adds code after referenced code
func (tr *Transformer) AppendAfter(reference string, source ...string) error {

	content, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lineNumber := -1
	for num, line := range lines {
		if !strings.Contains(line, reference) {
			continue
		}

		lineNumber = num + 1
		break
	}

	if lineNumber != -1 {
		lines = append(lines[:lineNumber], append(source, lines[lineNumber:]...)...)
	}

	fileContent := strings.Join(lines, "\n")
	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//AppendBefore adds code before passed reference
func (tr *Transformer) AppendBefore(reference string, source ...string) error {
	content, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lineNumber := -1
	for num, line := range lines {
		if !strings.Contains(line, reference) {
			continue
		}

		lineNumber = num + 1
		break
	}

	if lineNumber != -1 {
		lines = append(lines[:lineNumber-1], append(source, lines[lineNumber-1:]...)...)
	}

	fileContent := strings.Join(lines, "\n")
	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//RemoveLine removes a line starting with some passed code
func (tr *Transformer) RemoveLine(starting string) error {
	src, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(src), "\n")
	lineNum := -1
	for index, line := range lines {
		if strings.HasPrefix(line, starting) {
			lineNum = index
			break
		}
	}

	if lineNum > 0 {
		return nil
	}

	c := append(lines[:lineNum], lines[lineNum:]...)
	fileContent := strings.Join(c, "\n")

	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//RemoveBlock removes a block starting with passed expression
func (tr *Transformer) RemoveBlock(starting string) error {
	start, end, err := tr.FindBlockFor(starting)
	if err != nil {
		return err
	}

	if end < 0 {
		logrus.Warnf("could not find block in %v", tr.filePath)
		return nil
	}

	src, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(src), "\n")
	c := append(lines[:start-1], lines[end:]...)
	fileContent := strings.Join(c, "\n")

	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//SetBlockBody replaces block body starting with string
func (tr *Transformer) SetBlockBody(starting string, content ...string) error {
	start, end, err := tr.FindBlockFor(starting)
	if err != nil {
		return err
	}

	if end < 0 {
		return fmt.Errorf("could not find desired block in the file (%v)", tr.filePath)
	}

	src, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(src), "\n")
	c := append(lines[:start], append(content, lines[end-1:]...)...)
	fileContent := strings.Join(c, "\n")

	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//AppendToBlock adds source before block ends
func (tr *Transformer) AppendToBlock(startingExpr string, content ...string) error {
	_, end, err := tr.FindBlockFor(startingExpr)
	if err != nil {
		return err
	}

	src, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(src), "\n")
	c := append(lines[:end-1], append(content, lines[end-1:]...)...)
	fileContent := strings.Join(c, "\n")

	err = ioutil.WriteFile(tr.filePath, []byte(fileContent), 0755)
	return err
}

//FindBlockFor finds a block line start and end
func (tr *Transformer) FindBlockFor(startingExpr string) (int, int, error) {
	end, start := -1, -1

	src, err := ioutil.ReadFile(tr.filePath)
	if err != nil {
		return start, end, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, tr.filePath, string(src), 0)
	if err != nil {
		return start, end, err
	}

	lines := strings.Split(string(src), "\n")
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.StructType:
			line := fset.Position(x.Pos()).Line
			structDeclaration := fmt.Sprintf("%s\n", lines[line-1])

			if strings.Contains(structDeclaration, startingExpr) {
				start = line
				end = fset.Position(x.End()).Line
				return false
			}

		case *ast.BlockStmt:
			line := fset.Position(x.Lbrace).Line
			blockDeclaration := fmt.Sprintf("%s\n", lines[line-1])

			if strings.Contains(blockDeclaration, startingExpr) {
				start = line
				end = fset.Position(x.Rbrace).Line
			}

		}
		return true
	})

	return start, end, nil
}
