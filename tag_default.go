package envy

import (
	"context"
	"reflect"
)

const default_tagname = "default"

func WithDefaultTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		t.Default = string(field.Tag.Get(default_tagname))
		t.Value = t.Default
		return next.UnmarshalField(ctx, field)
	})
}
