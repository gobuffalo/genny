package genny

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

func (r *Suite) Test_WithLogger() {
	g := Background()

	g = WithRunner(g, func(gg Generator) error {
		gg.Logger().Infof("hi!")
		return nil
	})

	bb := &bytes.Buffer{}
	l := logrus.New()
	l.Out = bb

	g = WithLogger(g, l)

	r.NoError(g.Run())
	r.Contains(bb.String(), "hi!")
}
