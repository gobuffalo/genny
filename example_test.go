package genny_test

import (
	"context"
	"fmt"
	"go/build"
	"log"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
)

// do just cleans up variable log content
// such as GOPATH, step names, etc....
// without this Go Example tests won't work.
func do(fn func(r *genny.Runner)) {
	run := genny.DryRunner(context.Background())
	l := gentest.NewLogger()
	l.CloseFn = func() error {
		s := l.Stream.String()
		c := build.Default
		for _, src := range c.SrcDirs() {
			s = strings.Replace(s, src, "/go/src", -1)
		}
		s = strings.Replace(s, "\\", "/", -1)

		for i, line := range strings.Split(s, "\n") {
			if strings.Contains(line, "Step:") {
				s = strings.Replace(s, line, fmt.Sprintf("[DEBU] Step: %d", i+1), 1)
			}
		}
		fmt.Print(s)
		return nil
	}

	run.Logger = l
	fn(run)

}

func ExampleGenerator() {
	do(func(r *genny.Runner) {
		// create a new `*genny.Generator`
		g := genny.New()

		// add a file named `index.html` that has a body of `Hello\n`
		// to the generator
		g.File(genny.NewFileS("index.html", "Hello\n"))

		// add the generator to the `*genny.Runner`.
		r.With(g)

		// run the runner
		if err := r.Run(); err != nil {
			log.Fatal(err)
		}
	})
	// Output:
	// [DEBU] Step: 1
	// [DEBU] Chdir: /go/src/github.com/gobuffalo/genny
	// [DEBU] File: /go/src/github.com/gobuffalo/genny/index.html
}
