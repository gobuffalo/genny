package genny

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"
)

func testLogger(r *Runner) *bytes.Buffer {
	bb := &bytes.Buffer{}
	l := logrus.New()
	l.Out = bb
	l.SetLevel(logrus.DebugLevel)
	r.Logger = l
	return bb
}
func testRunner() (*Runner, *bytes.Buffer) {
	r := DryRunner(context.Background())
	return r, testLogger(r)
}
