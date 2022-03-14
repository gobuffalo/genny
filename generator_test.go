package genny

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2/internal/testdata"
	"github.com/gobuffalo/packd"
	"github.com/stretchr/testify/require"
)

var fixtures = func() packd.Box {
	box := packd.NewMemoryBox()
	box.AddString("foo.txt", "foo!")
	box.AddString("bar/baz.txt", "baz!")
	return box
}()

func Test_Generator_File(t *testing.T) {
	r := require.New(t)

	g := New()
	g.File(NewFile("foo.txt", strings.NewReader("hello")))

	run := DryRunner(context.Background())
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("foo.txt", f.Name())
	r.Equal("hello", f.String())
}

func Test_Generator_Box(t *testing.T) {
	r := require.New(t)

	g := New()
	r.NoError(g.Box(fixtures))

	run := DryRunner(context.Background())
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("bar/baz.txt", f.Name())
	r.Equal("baz!", f.String())

	f = res.Files[1]
	r.Equal("foo.txt", f.Name())
	r.Equal("foo!", f.String())
}

func Test_Generator_FS(t *testing.T) {
	r := require.New(t)

	g := New()
	r.NoError(g.FS(testdata.Data()))

	run := DryRunner(context.Background())
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 8)

	f := res.Files[0]
	r.Equal("bar/baz.txt", f.Name())
	r.Equal("baz!", f.String())

	f = res.Files[1]
	r.Equal("foo.txt", f.Name())
	r.Equal("foo!", f.String())
}

func Test_Generator_OnlyFS(t *testing.T) {
	td := []struct {
		name          string
		includePrefix []string
		includeSuffix []string
		fs            []struct {
			name    string
			content string
		}
	}{
		{
			name:          "nil_nil",
			includePrefix: nil,
			includeSuffix: nil,
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "one_nil",
			includePrefix: []string{"sky/star"},
			includeSuffix: []string{},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
			},
		},
		{
			name:          "many_nil",
			includePrefix: []string{"sky/star", "sky/sun"},
			includeSuffix: []string{},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "nil_one",
			includePrefix: nil,
			includeSuffix: []string{"light.txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "one_one",
			includePrefix: []string{"sky/s"},
			includeSuffix: []string{"light"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/sun/light", content: "sun-light"},
			},
		},
		{
			name:          "many_one",
			includePrefix: []string{"sky/.moon", "sky/sun", "sky/comet"},
			includeSuffix: []string{"light"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/sun/light", content: "sun-light"},
			},
		},
		{
			name:          "nil_many",
			includePrefix: nil,
			includeSuffix: []string{"light", "t.txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "one_many",
			includePrefix: []string{"sky/.moon"},
			includeSuffix: []string{"light", "txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
			},
		},
		{
			name:          "many_many",
			includePrefix: []string{"sky/star", "sky/.moon"},
			includeSuffix: []string{"light", "txt", "exe"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
			},
		},
	}

	for _, tc := range td {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			g := New()
			r.NoError(g.OnlyFS(testdata.Data(), tc.includePrefix, tc.includeSuffix))

			run := DryRunner(context.Background())
			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()
			r.Len(res.Commands, 0)
			r.Len(res.Files, len(tc.fs))

			for no, fe := range tc.fs {
				f := res.Files[no]
				r.Equal(fe.name, f.Name())
				r.Equal(fe.content, f.String())
			}
		})
	}
}

func Test_Generator_ExceptFS(t *testing.T) {
	td := []struct {
		name          string
		excludePrefix []string
		excludeSuffix []string
		fs            []struct {
			name    string
			content string
		}
	}{
		{
			name:          "nil_nil",
			excludePrefix: nil,
			excludeSuffix: nil,
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "one_nil",
			excludePrefix: []string{"sky/star"},
			excludeSuffix: []string{},
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "many_nil",
			excludePrefix: []string{"sky/star", "sky/sun"},
			excludeSuffix: []string{},
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
			},
		},
		{
			name:          "nil_one",
			excludePrefix: nil,
			excludeSuffix: []string{"txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/sun/light", content: "sun-light"},
			},
		},
		{
			name:          "one_one",
			excludePrefix: []string{"sky/s"},
			excludeSuffix: []string{".txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/.moon/light.txt", content: "moon-light.txt"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/sun/light", content: "sun-light"},
			},
		},
		{
			name:          "many_one",
			excludePrefix: []string{"sky/.moon", "sky/sun", "___"},
			excludeSuffix: []string{".txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
			},
		},
		{
			name:          "nil_many",
			excludePrefix: nil,
			excludeSuffix: []string{"light", "txt"},
			fs: []struct {
				name    string
				content string
			}{},
		},
		{
			name:          "one_many",
			excludePrefix: []string{"sky/.moon"},
			excludeSuffix: []string{"light", "txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/star/light.txt", content: "star-light.txt"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
		{
			name:          "many_many",
			excludePrefix: []string{"sky/star", "sky/.moon"},
			excludeSuffix: []string{"_____", "txt"},
			fs: []struct {
				name    string
				content string
			}{
				{name: "bar/baz.txt", content: "baz!"},
				{name: "foo.txt", content: "foo!"},
				{name: "sky/.moon/light", content: "moon-light"},
				{name: "sky/star/light", content: "star-light"},
				{name: "sky/sun/light", content: "sun-light"},
				{name: "sky/sun/light.txt", content: "sun-light.txt"},
			},
		},
	}

	for _, tc := range td {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			g := New()
			r.NoError(g.ExceptFS(testdata.Data(), tc.excludePrefix, tc.excludeSuffix))

			run := DryRunner(context.Background())
			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()
			r.Len(res.Commands, 0)
			r.Len(res.Files, len(tc.fs))

			for no, fe := range tc.fs {
				f := res.Files[no]
				r.Equal(fe.name, f.Name())
				r.Equal(fe.content, f.String())
			}
		})
	}
}

func Test_Generator_SelectiveFS(t *testing.T) {
	r := require.New(t)

	includePrefix := []string{"sky"}
	includeSuffix := []string{"txt"}
	excludePrefix := []string{"sky/.moon"}
	excludeSuffix := []string{}

	g := New()
	r.NoError(g.SelectiveFS(testdata.Data(), includePrefix, includeSuffix, excludePrefix, excludeSuffix))

	run := DryRunner(context.Background())
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	td := []struct {
		no      int
		name    string
		content string
	}{
		{name: "sky/star/light.txt", content: "star-light.txt"},
		{name: "sky/sun/light.txt", content: "sun-light.txt"},
	}
	for no, te := range td {
		f := res.Files[no]
		r.Equal(te.name, f.Name())
		r.Equal(te.content, f.String())
	}
}

func Test_Command(t *testing.T) {
	r := require.New(t)

	g := New()
	g.Command(exec.Command("echo", "hello"))

	run := DryRunner(context.Background())
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 0)

	c := res.Commands[0]
	r.Equal("echo hello", strings.Join(c.Args, " "))
}

func Test_Merge(t *testing.T) {
	r := require.New(t)

	g1 := New()
	g1.Root = "one"
	g1.RunFn(func(r *Runner) error {
		return r.File(NewFileS("a.txt", "a"))
	})
	g1.RunFn(func(r *Runner) error {
		return r.File(NewFileS("b.txt", "b"))
	})
	g1.Transformer(NewTransformer("*", func(f File) (File, error) {
		return NewFileS(f.Name(), strings.ToUpper(f.String())), nil
	}))

	g2 := New()
	g2.Root = "two"
	g2.RunFn(func(r *Runner) error {
		return r.File(NewFileS("c.txt", "c"))
	})
	g2.Transformer(NewTransformer("*", func(f File) (File, error) {
		return NewFileS(f.Name(), f.String()+"g2"), nil
	}))

	g1.RunFn(func(r *Runner) error {
		return r.File(NewFileS("d.txt", "d"))
	})
	g1.Transformer(NewTransformer("*", func(f File) (File, error) {
		return NewFileS(f.Name(), f.String()+"g1"), nil
	}))

	g1.Merge(g2)

	run := DryRunner(context.Background())
	r.NoError(run.With(g1))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 4)
}
