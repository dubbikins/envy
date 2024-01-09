package envy

import (
	"context"
	"reflect"
)

type TagHandler interface {
	UnmarshalField(context.Context, reflect.StructField) error
}

type Middleware func(next TagHandler) TagHandler

type TagHandlerFunc func(ctx context.Context, field reflect.StructField) error

func (f TagHandlerFunc) UnmarshalField(ctx context.Context, field reflect.StructField) error {
	return f(ctx, field)
}
