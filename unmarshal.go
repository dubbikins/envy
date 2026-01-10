package envy

import (
	"context"

	"github.com/dubbikins/envy/v2/tag"
)

func Unmarshal(s any, walkFns ...tag.ChainableWalkFn) (err error) {
	return tag.Walk(tag.Chain(tag.UnmarshalText, walkFns...), s)
}

func UnmarshalContext(ctx context.Context, s any, walkFns ...tag.ChainableWalkFn) (err error) {
	return  tag.WalkContext(ctx, tag.Chain(tag.UnmarshalText, walkFns...), s)
}

type OptionsFunc[Options any] func(Options) error

func New[T any](options ...OptionsFunc[*T]) (*T, error) {
	var o *T = new(T)
	for _, option := range options {
		if err := option(o); err != nil {
			return o, err
		}
	}
	return o, nil
}

func FromEnvironment[T any](t *T) error {
	return Unmarshal(t)
}





