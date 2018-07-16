package genny

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Transformer(t *testing.T) {
	table := []struct {
		in  File
		out File
	}{
		{
			in:  NewFile("hi.txt.bar", strings.NewReader("hi")),
			out: NewFile("hello.txt", strings.NewReader("hi")),
		},
		{
			in:  NewFile("hi.txt", strings.NewReader("hi")),
			out: NewFile("hi.txt", strings.NewReader("hi")),
		},
	}

	tf := NewTransformer(".bar", func(f File) (File, error) {
		return NewFile("hello.txt.bar", f), nil
	})
	tf.StripExt = true

	for _, tt := range table {
		t.Run(tt.in.Name()+"|"+tt.out.Name(), func(st *testing.T) {
			r := require.New(st)
			f, err := tf.Transform(tt.in)
			r.NoError(err)
			r.Equal(tt.out.Name(), f.Name())

			ob, err := ioutil.ReadAll(tt.out)
			r.NoError(err)
			fb, err := ioutil.ReadAll(f)
			r.NoError(err)
			r.Equal(string(ob), string(fb))
		})
	}
}

func Test_Replace(t *testing.T) {
	r := require.New(t)

	table := []struct {
		in      string
		out     string
		search  string
		replace string
	}{
		{in: "foo/-dot-git-keep", out: "foo/.git-keep", search: "-dot-", replace: "."},
		{in: "foo/dot-git-keep", out: "foo/dot-git-keep", search: "-dot-", replace: "."},
	}

	for _, tt := range table {
		in := NewFile(tt.in, nil)
		out, err := Replace(tt.search, tt.replace).Transform(in)
		r.NoError(err)
		r.Equal(tt.out, out.Name())
	}
}
