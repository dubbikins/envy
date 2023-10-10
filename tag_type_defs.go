package envy

import (
	"context"
	"reflect"
)

type TagHandlerFunc func(context.Context, reflect.StructField) error

func (f TagHandlerFunc) UnmarshalField(ctx context.Context, field reflect.StructField) error {
	return f(ctx, field)
}

type Middleware func(next TagHandler) TagHandler

type TagHandler interface {
	UnmarshalField(context.Context, reflect.StructField) error
}
type TagMiddleware interface {
	Pop() Middleware
	Push(...Middleware)
	Contents() string
	GetState() map[string]interface{}
	GetStateValue(key string) interface{}
}

type TextUnmarshallable interface {
	UnmarshalText(text []byte) error
}

type UnmarshallableFunc func(text []byte) error

func (f UnmarshallableFunc) UnmarshalText(text []byte) error { return f(text) }
