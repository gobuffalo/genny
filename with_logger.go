package genny

type withLogger struct {
	Generator
	logger Logger
}

func (w withLogger) Parent() Generator {
	return w.Generator
}

func (w withLogger) Logger() Logger {
	return w.logger
}

func (w *withLogger) setLogger(l Logger) {
	w.logger = l
}

// WithLogger will apply a new Logger to the entire tree
// of generators.
func WithLogger(g Generator, l Logger) Generator {
	if sl, ok := g.(setLogable); ok {
		sl.setLogger(l)
		return g
	}
	p := g.Parent()
	for p != nil {
		if sl, ok := p.(setLogable); ok {
			sl.setLogger(l)
			return g
		}
		p = p.Parent()
	}
	return g
}
