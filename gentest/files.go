package gentest

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/genny"
)

// CompareFiles compares a slice of expected filenames to a slice of
// genny.Files. Expected files can be listed in any order. An excellent choice
// for the actual files can be found in genny#Generator.Results under the Files
// attribute
func CompareFiles(exp []string, act []genny.File) error {
	if len(exp) != len(act) {
		return fmt.Errorf("len(exp) != len(act) [%d != %d]", len(exp), len(act))
	}

	var acts []string
	for _, f := range act {
		acts = append(acts, f.Name())
	}
	sort.Strings(exp)
	sort.Strings(acts)

	for i, n := range exp {
		if n != acts[i] {
			return fmt.Errorf("expected %q to match %q", acts, exp)
		}
	}
	return nil
}
