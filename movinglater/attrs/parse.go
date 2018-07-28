package attrs

import (
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/pkg/errors"
)

func Parse(arg string) (Attr, error) {
	arg = strings.TrimSpace(arg)
	attr := Attr{
		Original:   arg,
		commonType: "string",
	}
	if len(arg) == 0 {
		return attr, errors.New("argument can not be blank")
	}

	parts := strings.Split(arg, ":")
	attr.Name = name.New(parts[0])
	if len(parts) > 1 {
		attr.commonType = parts[1]
	}

	if len(parts) > 2 {
		attr.goType = parts[2]
	}

	return attr, nil
}

func ParseArgs(args ...string) (Attrs, error) {
	var attrs Attrs

	for _, arg := range args {
		a, err := Parse(arg)
		if err != nil {
			return attrs, errors.WithStack(err)
		}
		attrs = append(attrs, a)
	}

	return attrs, nil
}
