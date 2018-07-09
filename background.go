package genny

import (
	"context"
	"os/exec"
)

type background struct{}

var _ Generator = background{}

func (b background) Run() error {
	return nil
}

func (b background) Parent() Generator {
	return nil
}

func (b background) Context() context.Context {
	return nil
}

func (b background) Logger() Logger {
	return nil
}

func (b background) Cmd() *exec.Cmd {
	return nil
}

func (b background) File() File {
	return nil
}

// Background returns a new "blank" generator with context.Background()
func Background() Generator {
	return New(context.Background())
}
