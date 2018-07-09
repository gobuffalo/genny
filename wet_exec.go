package genny

import (
	"os"
	"os/exec"
)

// WetExec returns a new Generator that will execute commnds
// and wrap them with the STD in/out/err mappings
func WetExec(g Generator) Generator {
	return WithExec(g, func(cmd *exec.Cmd) error {
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	})
}
