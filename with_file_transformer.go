package genny

type FileTransformer interface {
	Transform(File) (File, error)
}

type FileTransformerFn func(File) (File, error)

type withFileTransformer struct {
	Generator
	fn FileTransformerFn
}

func (w withFileTransformer) Parent() Generator {
	return w.Generator
}

func (w withFileTransformer) String() string {
	return "genny.WithFileTransformer"
}

func (w withFileTransformer) Transform(f File) (File, error) {
	if w.fn == nil {
		return f, nil
	}
	return w.fn(f)
}

func WithFileTransformer(g Generator, fn FileTransformerFn) Generator {
	g = withFileTransformer{
		Generator: g,
		fn:        fn,
	}
	return g
}
