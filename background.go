package genny

import (
	"context"
	"os/exec"
)

type background struct{}

var _ Generator = background{}
var _ Execer = background{}
var _ Commandable = background{}
var _ Fileable = background{}
var _ FileHandler = background{}
var _ FileTransformer = background{}

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

func (b background) String() string {
	return "genny.background"
}

func (b background) Exec(*exec.Cmd) error {
	return nil
}

func (b background) Handle(File) error {
	return nil
}

func (b background) Transform(f File) (File, error) {
	return f, nil
}

// Background returns a new "blank" generator with context.Background()
func Background() Generator {
	return New(context.Background())
}
