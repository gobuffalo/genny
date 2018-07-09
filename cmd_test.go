package genny

func (r *Suite) Test_Cmds() {
	cmd := r.Command("bad-command", "$$$;;;")
	g := Background()
	g = WithCmd(g, cmd)
	g = DryExec(g)
	cmds := Cmds(g)
	r.Len(cmds, 1)
}

func (r *Suite) Test_Cmds_FiFo() {
	g := Background()

	c1 := r.Command("c1")
	g = WithCmd(g, c1)

	c2 := r.Command("c2")
	g = WithCmd(g, c2)

	g = DryExec(g)

	cmds := Cmds(g)
	r.Len(cmds, 2)

	r.Equal("c1", cmds[0].Args[0])
	r.Equal("c2", cmds[1].Args[0])
}

func (r *Suite) Test_Cmds_PreventDups_ExactCmd() {
	cmd := r.Command("bad-command", "$$$;;;")
	g := Background()
	g = WithCmd(g, cmd)
	g = WithCmd(g, cmd)
	g = DryExec(g)
	cmds := Cmds(g)
	r.Len(cmds, 1)
}

func (r *Suite) Test_Cmds_AllowsDups_DifferentCmd() {
	cmd := r.Command("bad-command", "$$$;;;")
	g := Background()
	g = WithCmd(g, cmd)
	g = WithCmd(g, r.Command("bad-command", "$$$;;;"))
	g = DryExec(g)
	cmds := Cmds(g)
	r.Len(cmds, 2)
}
