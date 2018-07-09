package genny

// RunFn is a type representing the Run function of a Generator
type RunFn func() error

// GRunFn is used by `WithRunner`
type GRunFn func(Generator) error

type withRunner struct {
	Generator
	run GRunFn
}

func (w withRunner) Parent() Generator {
	return w.Generator
}

func (w withRunner) Run() error {
	return w.run(w.Generator)
}

func (w withRunner) String() string {
	return "genny.WithRunner"
}

// WithRunner adds an arbitrary Run call to the generator
// stack
func WithRunner(g Generator, rf GRunFn) Generator {
	g = withRunner{
		Generator: g,
		run:       rf,
	}
	return g
}
