package genny

import (
	"strings"
)

func (r *Suite) Test_DryRun() {
	g := Background()
	g, bb := withTestLogger(g)
	g = DryRun(g)
	g = WithCmd(g, r.Command("foo", "bar"))
	g = WithFile(g, NewFile("baz.txt", strings.NewReader("bazzy")))
	r.NoError(Run(g))
	r.Contains(bb.String(), "[genny.WithFilesHandler] baz.txt")
	r.Contains(bb.String(), "[genny.WithExec] foo bar")
}
