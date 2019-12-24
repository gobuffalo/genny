package gomods

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/here"
)

func New(name string, path string) (*genny.Group, error) {
	g := &genny.Group{}

	init, err := Init(name, path)
	if err != nil {
		return g, err
	}
	g.Add(init)

	tidy, err := Tidy(path, false)
	if err != nil {
		return g, err
	}
	g.Add(tidy)
	return g, nil
}

func Init(name string, root string) (*genny.Generator, error) {
	if len(root) == 0 || root == "." {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		root = pwd
	}

	if len(name) == 0 {
		name = path.Base(root)
		if info, err := here.Dir(root); err == nil {
			name = info.ImportPath
		}
	}

	name = strings.Replace(name, "\\", "/", -1)
	name = strings.TrimPrefix(name, "/")

	g := genny.New()
	g.StepName = "go:mod:init:" + name
	g.RunFn(func(r *genny.Runner) error {
		return r.Chdir(root, func() error {
			args := []string{"mod", "init"}
			if len(name) > 0 {
				args = append(args, name)
			}
			cmd := exec.Command(genny.GoBin(), args...)
			return r.Exec(cmd)
		})
	})
	return g, nil
}
