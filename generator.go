package genny

import (
	"io/fs"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/packd"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generator is the basic type for generators to use
type Generator struct {
	StepName     string
	Should       func(*Runner) bool
	Root         string
	ErrorFn      func(error)
	runners      []RunFn
	transformers []Transformer
	moot         *sync.RWMutex
}

// New, well-formed, generator
func New() *Generator {
	g := &Generator{
		StepName:     stepName(),
		runners:      []RunFn{},
		moot:         &sync.RWMutex{},
		transformers: []Transformer{},
	}
	return g
}

// File adds a file to be run when the generator is run
func (g *Generator) File(f File) {
	g.RunFn(func(r *Runner) error {
		return r.File(f)
	})
}

func (g *Generator) Transform(f File) (File, error) {
	g.moot.RLock()
	defer g.moot.RUnlock()
	var err error
	for _, t := range g.transformers {
		f, err = t.Transform(f)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

// Transformer adds a file transform to the generator
func (g *Generator) Transformer(t Transformer) {
	g.moot.Lock()
	defer g.moot.Unlock()
	g.transformers = append(g.transformers, t)
}

// Command adds a command to be run when the generator is run
func (g *Generator) Command(cmd *exec.Cmd) {
	g.RunFn(func(r *Runner) error {
		return r.Exec(cmd)
	})
}

// Box walks through a packr.Box and adds Files for each entry
// in the box.
func (g *Generator) Box(box packd.Walker) error {
	return box.Walk(func(path string, f packd.File) error {
		g.File(NewFile(path, f))
		return nil
	})
}

// ExceptFS walks through a given fs.FS and adds all files into Runner's
// virtual disk, except files that starts with excludePrefix and ends with
// excludeSuffix.
// includePrefix and includeSuffix will be interpreted as AND condition
// but an empty option will not be used. So if you give a prefix but
// suffix as nil, it will only check if files starts with the given prefix.
// It is usuful if you want to exclude entire subdirs.
func (g *Generator) ExceptFS(fsys fs.FS, excludePrefix []string, excludeSuffix []string) error {
	return g.SelectiveFS(fsys, nil, nil, excludePrefix, excludeSuffix)
}

// OnlyFS walks through a given fs.FS and adds only matching files into
// Runner's virtual disk, matching files that starts with includePrefix
// and ends with includeSuffix.
// includePrefix and includeSuffix will be interpreted as AND condition
// but an empty option will not be used. So if you give a prefix but
// suffix as nil, it will only check if files starts with the given prefix.
// It is usuful if you want to include specific subdirs only.
func (g *Generator) OnlyFS(fsys fs.FS, includePrefix []string, includeSuffix []string) error {
	return g.SelectiveFS(fsys, includePrefix, includeSuffix, nil, nil)
}

// FS walks through a fs.FS and adds Files for each entry.
func (g *Generator) FS(fsys fs.FS) error {
	return g.SelectiveFS(fsys, nil, nil, nil, nil)
}

// FS walks through a fs.FS and adds Files for each entry.
func (g *Generator) SelectiveFS(fsys fs.FS, includePrefix, includeSuffix, excludePrefix, excludeSuffix []string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		isIncluded := false
		isExcluded := false
		if len(includePrefix) > 0 {
			for _, fix := range includePrefix {
				if strings.HasPrefix(path, fix) {
					isIncluded = true
					break
				}
			}
		} else {
			isIncluded = true
		}
		inSuffix := false
		if len(includeSuffix) > 0 {
			for _, fix := range includeSuffix {
				if strings.HasSuffix(path, fix) {
					inSuffix = true
					break
				}
			}
		} else {
			inSuffix = true
		}
		isIncluded = isIncluded && inSuffix

		for _, fix := range excludePrefix {
			if strings.HasPrefix(path, fix) {
				isExcluded = true
				break
			}
		}
		exSuffix := false
		if len(excludeSuffix) == 0 {
			exSuffix = isExcluded
		}
		for _, fix := range excludeSuffix {
			if len(excludePrefix) == 0 {
				isExcluded = true
			}
			exSuffix = isExcluded && strings.HasSuffix(path, fix)
			if exSuffix {
				break
			}
		}
		isExcluded = exSuffix

		if !isIncluded || isExcluded {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		g.File(NewFile(path, f))
		return nil
	})
}

// RunFn adds a generic "runner" function to the generator.
func (g *Generator) RunFn(fn RunFn) {
	g.moot.Lock()
	defer g.moot.Unlock()
	g.runners = append(g.runners, fn)
}

func (g1 *Generator) Merge(g2 *Generator) {
	g2.moot.Lock()
	g1.moot.Lock()
	g1.runners = append(g1.runners, g2.runners...)
	g1.transformers = append(g1.transformers, g2.transformers...)
	g1.moot.Unlock()
	g2.moot.Unlock()
}
