package genny

import (
	"bytes"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func (r *Suite) Test_WithExec() {
	g := Background()
	g = WithCmd(g, r.Command("foo", "bar", "-v"))
	g = DryExec(g)
	g = WithExec(g, func(cmd *exec.Cmd) error {
		return cmd.Run()
	})

	bb := &bytes.Buffer{}
	l := logrus.New()
	l.Out = bb
	g = WithLogger(g, l)
	err := g.Run()
	r.Error(err)

	r.Contains(bb.String(), "[genny.WithExec] foo bar -v")
}
