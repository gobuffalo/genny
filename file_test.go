package genny

import (
	"io/ioutil"
	"strings"
)

func (r *Suite) Test_NewFile() {
	input := "hi"
	f := NewFile("foo.txt", strings.NewReader(input))
	r.NotNil(f)
	r.Equal("foo.txt", f.Name())
	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal(input, string(b))
}
