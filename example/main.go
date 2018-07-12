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

	attrs, err := attrs.ParseArgs(args...)
	if err != nil {
		log.Fatal(err)
	}

	err = r.WithFn(func() (*genny.Generator, error) {
		return pop.Model(attrs)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = r.WithFn(func() (*genny.Generator, error) {
		return fizz.FizzMigration(attrs)
	})
	if err != nil {
		log.Fatal(err)
	}

	r.With(MyCustomGenerator(attrs))

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

func MyCustomGenerator(ats attrs.Attrs) *genny.Generator {
	g := genny.New()
	g.File(genny.NewFile("templates/index.plush.html", strings.NewReader(template)))
	g.Command(exec.Command("asdf", "asdfasdf"))

	ctx := plush.NewContext()
	ctx.Set("model", ats[0].Name)
	g.Transformer(plushgen.Transformer(ctx))
	return g
}

const template = `
<h1>Hello <%= model.Model() %></h1>
`
