package genny

// Logger interface for a logger to be used
// with genny. Logrus is 100% compatible.
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type setLogable interface {
	setLogger(Logger)
}
