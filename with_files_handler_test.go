package genny

import "strings"

func (r *Suite) Test_WithFilesHandler() {
	g := Background()

	var names []string
	g = WithFilesHandler(g, func(f File) error {
		names = append(names, f.Name())
		return nil
	})

	g = WithFileFromReader(g, "foo.txt", strings.NewReader("foo"))
	g = WithFileFromReader(g, "bar.txt", strings.NewReader("bar"))

	r.NoError(Run(g))
	r.Len(names, 2)
}
