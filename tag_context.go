package envy

import (
	"context"
	"fmt"
)

type contextKey string

func getContextKey(key string) contextKey {
	return contextKey(key)
}
func WithTagContext(ctx context.Context, tag *tag) context.Context {
	return context.WithValue(ctx, getContextKey("_!_envy_tag_ctx_!_"), tag)
}
func GetTagContext(ctx context.Context) (*tag, error) {
	ectx, ok := ctx.Value(getContextKey("_!_envy_tag_ctx_!_")).(*tag)
	if !ok {
		return nil, fmt.Errorf("missing envy tag context")
	}
	return ectx, nil
}
func MustGetTagContext(ctx context.Context) *tag {
	ectx, ok := ctx.Value(getContextKey("_!_envy_tag_ctx_!_")).(*tag)
	if !ok {
		panic("missing envy tag context")
	}
	return ectx
}
