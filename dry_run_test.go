package genny

import "strings"

func (r *Suite) Test_DryRun() {
	g := Background()
	g = WithCmd(g, r.Command("foo", "bar"))
	g = WithFile(g, NewFile("baz.txt", strings.NewReader("bazzy")))
	g = DryRun(g)
	r.NoError(g.Run())
}
