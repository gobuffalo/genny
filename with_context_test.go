package genny

import (
	"context"
	"time"
)

func (r *Suite) Test_Context() {
	g := Background()
	r.Equal(g.Context(), context.Background())

	ctx, cancel := context.WithTimeout(g.Context(), time.Second)
	defer cancel()

	g = WithContext(g, ctx)
	r.NotEqual(g.Context(), context.Background())
	r.Equal(g.Context(), ctx)
}
