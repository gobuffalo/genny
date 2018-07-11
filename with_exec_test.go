package genny

func (r *Suite) Test_WithExec() {
	g := Background()
	g, bb := withTestLogger(g)
	g = DryExec(g)
	g = WithCmd(g, r.Command("foo", "bar", "-v"))

	err := g.Run()
	r.NoError(err)

	r.Contains(bb.String(), "[genny.WithExec] foo bar -v")
}
