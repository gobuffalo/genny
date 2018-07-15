package attrs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseNamedArgs(t *testing.T) {
	r := require.New(t)

	na, err := ParseNamedArgs("widget", "name", "age:int")
	r.NoError(err)

	r.Equal("widget", na.Name.String())
	r.Len(na.Attrs, 2)
}

func Test_ParseNamedArgs_NoName(t *testing.T) {
	r := require.New(t)

	_, err := ParseNamedArgs()
	r.Error(err)

}
