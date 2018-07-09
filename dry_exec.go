package genny

import (
	"os/exec"
)

// DryExec will just log out the command and it's arguments
// it will NOT run the command
func DryExec(g Generator) Generator {
	g = WithExec(g, func(cmd *exec.Cmd) error {
		return nil
	})
	return g
}
