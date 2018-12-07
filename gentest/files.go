package gentest

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/genny"
)

// CompareFiles ...
func CompareFiles(exp []string, act []genny.File) error {
	if len(exp) != len(act) {
		return fmt.Errorf("len(exp) != len(act) [%d != %d]", len(exp), len(act))
	}
	sort.Strings(exp)
	sort.Slice(act, func(a, b int) bool {
		return act[a].Name() < act[b].Name()
	})

	for i, f := range act {
		e := exp[i]
		a := f.Name()
		if a != e {
			return fmt.Errorf("expect %q got %q", a, e)
		}
	}
	return nil
}
