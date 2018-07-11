package genny

func (r *Suite) Test_WithLogger() {
	g := Background()
	g, bb := withTestLogger(g)
	g = WithRunner(g, func(gg Generator) error {
		gg.Logger().Infof("hi!")
		return nil
	})

	r.NoError(g.Run())
	r.Contains(bb.String(), "hi!")
}
