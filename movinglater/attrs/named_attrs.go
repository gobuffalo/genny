package attrs

import (
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/pkg/errors"
)

type NamedAttrs struct {
	Name  name.Ident
	Attrs Attrs
}

func (n NamedAttrs) String() string {
	x := []string{n.Name.Original}
	for _, a := range n.Attrs {
		x = append(x, a.String())
	}
	return strings.Join(x, " ")
}

func ParseNamedArgs(args ...string) (NamedAttrs, error) {
	var na NamedAttrs
	if len(args) < 1 {
		return na, errors.New("requires a name argument")
	}
	na.Name = name.New(args[0])
	if len(args) > 1 {
		var err error
		if na.Attrs, err = ParseArgs(args[1:]...); err != nil {
			return na, err
		}
	}
	return na, nil
}
