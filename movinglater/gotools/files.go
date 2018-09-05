package gotools

import (
	"bytes"
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

func GoFiles(dir string) ([]string, error) {
	var files []string

	pwd, err := os.Getwd()
	if err != nil {
		return files, errors.WithStack(err)
	}
	if dir == "" {
		dir = pwd
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		path = strings.TrimPrefix(path, pwd+"/")
		if strings.Contains(path, ".git") || strings.Contains(path, "node_modules") || strings.Contains(path, "vendor"+string(os.PathSeparator)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	return files, nil
}

type ParsedFile struct {
	File    genny.File
	FileSet *token.FileSet
	Ast     *ast.File
	Lines   []string
}

func parseFile(gf genny.File) (ParsedFile, error) {
	name := gf.Name()
	pf := ParsedFile{
		FileSet: token.NewFileSet(),
	}

	gf, err := beforeParse(gf)
	if err != nil {
		return pf, errors.WithStack(err)
	}

	src := gf.String()
	f, err := parser.ParseFile(pf.FileSet, gf.Name(), src, 0)
	if err != nil {
		return pf, errors.WithStack(err)
	}
	pf.Ast = f

	pf.Lines = strings.Split(src, "\n")
	pf.File = genny.NewFile(name, gf)
	return pf, nil
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
