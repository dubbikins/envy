package envy

import (
	"context"
	"reflect"
)

type Middleware func(next TagHandler) TagHandler

type TagHandler interface {
	UnmarshalField(context.Context, reflect.StructField) error
}
type TagMiddleware interface {
	Pop() Middleware
	Push(...Middleware)
	Reset()
}

type TagParserFunc func(ctx context.Context, field reflect.StructField) error

type TagHandlerFunc TagParserFunc

func (f TagHandlerFunc) UnmarshalField(ctx context.Context, field reflect.StructField) error {
	return f(ctx, field)
}
