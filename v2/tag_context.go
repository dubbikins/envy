package envy

import (
	"context"
	"fmt"
)

type contextKey string
type ContextKey contextKey

func getContextKey(key string) contextKey {
	return contextKey(key)
}
func WithTagContext(ctx context.Context, tag *Tag) context.Context {
	return context.WithValue(ctx, getContextKey("_!_envy_tag_ctx_!_"), tag)
}
func GetTagContext(ctx context.Context) (*Tag, error) {
	ectx, ok := ctx.Value(getContextKey("_!_envy_tag_ctx_!_")).(*Tag)
	if !ok {
		return nil, fmt.Errorf("missing envy tag in context")
	}
	return ectx, nil
}

func WithOptionsContext(ctx context.Context, options *Options) context.Context {
	return context.WithValue(ctx, getContextKey("_!_envy_options_ctx_!_"), options)
}
func GetOptionsContext(ctx context.Context) (*Options, error) {
	ectx, ok := ctx.Value(getContextKey("_!_envy_options_ctx_!_")).(*Options)
	if !ok {
		return nil, fmt.Errorf("missing envy options in context")
	}
	return ectx, nil
}
func MustGetTagContext(ctx context.Context) *Tag {
	ectx, ok := ctx.Value(getContextKey("_!_envy_tag_ctx_!_")).(*Tag)
	if !ok {
		panic("missing envy tag in context")
	}
	return ectx
}

func GetShouldSkip(ctx context.Context) (*Options, error) {
	ectx, ok := ctx.Value(getContextKey("_!_envy_options_ctx_!_")).(*Options)
	if !ok {
		return nil, fmt.Errorf("missing envy options in context")
	}
	return ectx, nil
}
