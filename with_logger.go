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

func (w withLogger) String() string {
	return "genny.WithLogger"
}

// WithLogger will apply a new Logger to the entire tree
// of generators.
func WithLogger(g Generator, l Logger) Generator {
	g = withLogger{
		Generator: g,
		logger:    l,
	}
	return g
}
