package genny

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewFile(t *testing.T) {
	r := require.New(t)

	input := "hi"
	f := NewFile("foo.txt", strings.NewReader(input))
	r.NotNil(f)
	r.Equal("foo.txt", f.Name())
	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal(input, string(b))
}
