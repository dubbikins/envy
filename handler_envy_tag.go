package envy

import (
	"context"
	"reflect"
)

const envy_global_tagname = "envy"

func WithEnvyGlobalTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) (err error) {
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		t.Name = field.Tag.Get(envy_global_tagname)
		if t.Name == "-" {
			t.Skip = true
			return nil
		}
		return next.UnmarshalField(ctx, field)
	})
}
