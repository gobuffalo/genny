package genny

import (
	"bytes"
	"os"

	"github.com/sirupsen/logrus"
)

func (r *Suite) Test_WetExec() {
	g := Background()

	g = WithCmd(g, r.Command("bad-command", "bad-args"))
	g = WetExec(g)

	bb := &bytes.Buffer{}
	l := logrus.New()
	l.Out = bb
	g = WithLogger(g, l)

	err := g.Run()
	r.Error(err)

	r.Contains(err.Error(), "executable file not found in")
}

func (r *Suite) Test_WetExec_SetsIO() {
	g := Background()

	cmd := r.Command("bad-command", "bad args")
	g = WithCmd(g, cmd)

	r.NotEqual(cmd.Stdin, os.Stdin)
	r.NotEqual(cmd.Stderr, os.Stderr)
	r.NotEqual(cmd.Stdout, os.Stdout)

	g = WetExec(g)
	r.Error(g.Run())

	r.Equal(cmd.Stdin, os.Stdin)
	r.Equal(cmd.Stderr, os.Stderr)
	r.Equal(cmd.Stdout, os.Stdout)
}
