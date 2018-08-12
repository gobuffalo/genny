package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/example/fizz"
	"github.com/gobuffalo/genny/example/pop"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/plush"
)

func main() {
	r := genny.DryRunner(context.Background())
	r.Root = "./coke"

	args := []string{"widget", "name", "birthdate:timstamp"}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	ats, err := attrs.ParseNamedArgs(args...)
	if err != nil {
		log.Fatal(err)
	}

	err = r.WithFn(func() (*genny.Generator, error) {
		return pop.Model(ats)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = r.WithFn(func() (*genny.Generator, error) {
		return fizz.FizzMigration(ats)
	})
	if err != nil {
		log.Fatal(err)
	}

	r.With(MyCustomGenerator(ats))

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

func MyCustomGenerator(ats attrs.NamedAttrs) *genny.Generator {
	g := genny.New()
	g.File(genny.NewFile("templates/index.plush.html", strings.NewReader(template)))
	g.Command(exec.Command("asdf", "asdfasdf"))

	ctx := plush.NewContext()
	ctx.Set("model", ats)
	g.Transformer(plushgen.Transformer(ctx))
	return g
}

const template = `
<h1>Hello <%= model.Name.Model() %></h1>
`
