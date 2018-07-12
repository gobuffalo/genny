package attrs

import (
	"testing"

	"github.com/markbates/inflect"
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

func Test_Attr_GoType(t *testing.T) {

	tt := []struct {
		ct string
		gt string
	}{
		{"timestamp", "time.Time"},
		{"datetime", "time.Time"},
		{"date", "time.Time"},
		{"time", "time.Time"},
	}

	for _, test := range tt {
		t.Run(test.ct+"/"+test.gt, func(st *testing.T) {
			r := require.New(st)
			a := Attr{commonType: inflect.Name(test.ct)}
			r.Equal(test.gt, a.GoType().String())
		})
	}

}
