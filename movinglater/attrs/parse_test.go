package attrs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	attrs := Attrs{
		{Original: "foo", goType: "string", commonType: "string", Name: "foo"},
		{Original: "foo:int", goType: "int", commonType: "int", Name: "foo"},
		{Original: "foo:timestamp", goType: "time.Time", commonType: "timestamp", Name: "foo"},
		{Original: "foo:text:exec.Command", goType: "exec.Command", commonType: "text", Name: "foo"},
	}

	for _, a := range attrs {
		t.Run(a.Original, func(st *testing.T) {
			r := require.New(st)
			attr, err := Parse(a.Original)
			r.NoError(err)
			r.Equal(a.Original, attr.Original)
			r.Equal(a.goType, attr.GoType())
			r.Equal(a.commonType, attr.CommonType())
			r.Equal(a.Name, attr.Name)
		})
	}
}
