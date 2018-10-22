package genny

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewStep(t *testing.T) {
	r := require.New(t)

	_, err := NewStep(nil, 0)
	r.Error(err)

	s, err := NewStep(New(), 0)
	r.NoError(err)
	r.NotZero(s)
}
