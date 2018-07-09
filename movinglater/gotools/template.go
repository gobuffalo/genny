package gotools

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// WithTemplate returns a generator who's previous file has been run through text/template
func WithTemplate(g genny.Generator, data interface{}, helpers template.FuncMap) (genny.Generator, error) {
	fa, ok := g.(genny.Fileable)
	if !ok {
		return g, genny.ErrNilFile
	}
	f, err := renderWithTemplate(fa.File(), data, helpers)
	if err != nil {
		return g, errors.WithStack(err)
	}
	return genny.WithFile(g, f), nil
}

func renderWithTemplate(f genny.File, data interface{}, helpers template.FuncMap) (genny.File, error) {
	if f == nil {
		return f, genny.ErrNilFile
	}
	path := f.Name()
	t := template.New(path)
	if helpers != nil {
		t = t.Funcs(helpers)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return f, errors.WithStack(err)
	}
	t, err = t.Parse(string(b))
	if err != nil {
		return f, errors.WithStack(err)
	}

	var bb bytes.Buffer
	if err = t.Execute(&bb, data); err != nil {
		err = errors.WithStack(err)
		return f, errors.WithStack(err)
	}
	return genny.NewFile(path, &bb), nil
}
