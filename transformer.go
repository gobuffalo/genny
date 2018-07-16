package genny

import (
	"strings"

	"github.com/pkg/errors"
)

type TransformerFn func(File) (File, error)

type Transformer struct {
	Ext      string
	StripExt bool
	fn       TransformerFn
}

func (t Transformer) Transform(f File) (File, error) {
	if !HasExt(f, t.Ext) {
		return f, nil
	}
	if t.fn == nil {
		return f, nil
	}
	f, err := t.fn(f)
	if err != nil {
		return f, errors.WithStack(err)
	}
	if t.StripExt {
		return StripExt(f, t.Ext), nil
	}
	return f, nil
}

func NewTransformer(ext string, fn TransformerFn) Transformer {
	return Transformer{
		Ext: ext,
		fn:  fn,
	}
}

func Replace(search string, replace string) Transformer {
	return NewTransformer("*", func(f File) (File, error) {
		name := f.Name()
		name = strings.Replace(name, search, replace, -1)
		return NewFile(name, f), nil
	})
}
