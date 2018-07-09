package genny

import "strings"

func (r *Suite) Test_WithFilesHandler() {
	g := Background()
	g = WithFileFromReader(g, "foo.txt", strings.NewReader("foo"))
	g = WithFileFromReader(g, "bar.txt", strings.NewReader("bar"))

	var names []string
	g = WithFilesHandler(g, func(f File) error {
		names = append(names, f.Name())
		return nil
	})

	r.NoError(g.Run())
	r.Len(names, 2)
}
