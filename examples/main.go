package main

import (
	"log"
	"os/exec"
	"strings"

	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/genny/movinglater/packrgen"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

func main() {
	// setup a new generator:
	g := genny.Background()

	// wrap the previous generator with a command to run:
	g = genny.WithCmd(g, exec.CommandContext(g.Context(), "echo", "hello from the echo command!"))

	// wrap the previous generator with a file: (path to write to, reader to read from)
	g = genny.WithFileFromReader(g, "examples/output/foo.txt", strings.NewReader("Hello {{.}} <%= name %>"))

	// process the previously declared file using the text/template package:
	g, err := gotools.WithTemplate(g, "mark", template.FuncMap{})
	if err != nil {
		log.Fatal(err)
	}

	// process the previously declared file using the plush package:
	ctx := plush.NewContextWithContext(g.Context())
	ctx.Set("name", "Bates")
	g, err = plushgen.WithTemplate(g, ctx)
	if err != nil {
		log.Fatal(err)
	}

	// run another command
	g = genny.WithCmd(g, exec.CommandContext(g.Context(), "echo", "almost finished"))

	// add the contents of a packr box to the generators
	g, err = packrgen.WithBox(g, packr.NewBox("../"), nil)
	if err != nil {
		log.Fatal(err)
	}

	// add another file
	g = genny.WithFileFromReader(g, "examples/output/baz/bar.txt", strings.NewReader("plain text"))

	g = gotools.WithGoGet(g, "github.com/gobuffalo/envy", "-v")

	// wrap in a "dry runner" so files and commands are echoed to the screen, but not executed:
	g = genny.DryRun(g)

	// actually run the generators:
	err = g.Run()
	if err != nil {
		log.Fatal(err)
	}
}
