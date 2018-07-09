package genny

func (r *Suite) Test_WithCmd() {
	g := Background()

	cmd := r.Command("foo", "bar")
	g = WithCmd(g, cmd)

	ca, ok := g.(Commandable)
	r.True(ok)
	r.Equal(cmd, ca.Cmd())
}
