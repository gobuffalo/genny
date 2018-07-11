package genny

import (
	"os"
)

func (r *Suite) Test_WetExec() {
	g := Background()
	g = WetExec(g)

	g = WithCmd(g, r.Command("bad-command", "bad-args"))

	err := g.Run()
	r.Error(err)

	r.Contains(err.Error(), "executable file not found in")
}

func (r *Suite) Test_WetExec_SetsIO() {
	g := Background()
	g = WetExec(g)

	cmd := r.Command("bad-command", "bad args")
	g = WithCmd(g, cmd)

	r.NotEqual(cmd.Stdin, os.Stdin)
	r.NotEqual(cmd.Stderr, os.Stderr)
	r.NotEqual(cmd.Stdout, os.Stdout)

	r.Error(g.Run())

	r.Equal(cmd.Stdin, os.Stdin)
	r.Equal(cmd.Stderr, os.Stderr)
	r.Equal(cmd.Stdout, os.Stdout)
}
