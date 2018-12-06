package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/gobuffalo/genny"
)

func main() {
	run := genny.DryRunner(context.Background())

	g := genny.New()
	g.File(genny.NewFileS("index.html", "Hello\n"))
	g.Command(exec.Command("go", "env"))
	run.With(g)

	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
