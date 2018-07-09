<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://godoc.org/github.com/gobuffalo/genny"><img src="https://godoc.org/github.com/gobuffalo/genny?status.svg" alt="GoDoc" /></a>
<a href="https://travis-ci.org/gobuffalo/genny"><img src="https://travis-ci.org/gobuffalo/genny.svg?branch=master" alt="Build Status" /></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/genny"><img src="https://goreportcard.com/badge/github.com/gobuffalo/genny" alt="Go Report Card" /></a>
</p>

# Genny [EXPERIMENTAL]

**EXPERIMENTAL** - APIs can change without notice. You've been warned.

## What Is Genny?

Genny is a _framework_ for writing modular generators, it however, doesn't actually generate anything. It just makes it easier for you to. :)

## The `Generator` Interface

Genny was inspired by the [`context`](https://golang.org/pkg/context/) design and foregoes any configuration for a "wrapping" pattern. See the example at the bottom of the file or in the `./examples` folder.

```go
type Generator interface {
	Context() context.Context
	Run() error
	Parent() Generator
	Logger() Logger
}
```

## Documentation

For right now the [GoDoc](https://godoc.org/github.com/gobuffalo/genny) and the source/tests are best documentation as the APIs are currently in flux.

## Example

```go
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

	// wrap in a "dry runner" so files and commands are echoed to the screen, but not executed:
	g = genny.DryRun(g)

	// actually run the generators:
	err = g.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```
