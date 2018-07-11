package genny

func (r *Suite) Test_DryExec() {

	g := Background()
	g, bb := withTestLogger(g)
	g = DryExec(g)

	g = WithCmd(g, r.Command("bad-command", "$$$;;;"))

	r.NoError(g.Run())

	r.Contains(bb.String(), "[genny.WithExec] bad-command $$$;;;")
}
