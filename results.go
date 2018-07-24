package genny

import "os/exec"

type Results struct {
	Files    []File
	Commands []*exec.Cmd
}
