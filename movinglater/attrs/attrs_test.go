package attrs

import (
	"testing"

	"github.com/markbates/inflect"
	"github.com/stretchr/testify/require"
)

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
