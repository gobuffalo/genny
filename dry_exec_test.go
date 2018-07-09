package genny

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

func (r *Suite) Test_DryExec() {

	g := Background()
	g = WithCmd(g, r.Command("bad-command", "$$$;;;"))
	g = DryExec(g)

	bb := &bytes.Buffer{}
	l := logrus.New()
	l.Out = bb
	g = WithLogger(g, l)
	r.NoError(g.Run())

	r.Contains(bb.String(), "[genny.WithExec] bad-command $$$;;;")
}
