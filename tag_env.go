package envy

import (
	"context"
	"os"
	"reflect"
)

const env_tagname = "env"

func WithEnvTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		t.Name = field.Tag.Get(env_tagname)
		t.Value = os.Getenv(t.Name)
		return next.UnmarshalField(ctx, field)
	})
}
