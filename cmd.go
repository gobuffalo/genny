package genny

import (
	"fmt"
	"os/exec"
)

// Commandable interface for generators that return
// a command
type Commandable interface {
	Cmd() *exec.Cmd
}

// Cmds returns all of the commands up to the parent
// generator.
func Cmds(g Generator) []*exec.Cmd {
	var cmds []*exec.Cmd
	t := map[string]bool{}
	fp := func(c *exec.Cmd) {
		if c != nil {
			key := fmt.Sprint(c)
			if _, ok := t[key]; !ok {
				cmds = append([]*exec.Cmd{c}, cmds...)
				t[key] = true
			}
		}
	}
	if cm, ok := g.(Commandable); ok {
		fp(cm.Cmd())
	}

	p := g.Parent()
	for p != nil {
		if cm, ok := p.(Commandable); ok {
			fp(cm.Cmd())
		}
		p = p.Parent()
	}
	return cmds
}
