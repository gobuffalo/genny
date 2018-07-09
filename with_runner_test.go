package genny

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithRunner(t *testing.T) {
	r := require.New(t)

	var ran bool
	bg := WithRunner(Background(), func(Generator) error {
		ran = true
		return nil
	})
	r.NoError(bg.Run())
	r.True(ran)
}
