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
		if t.Name == "-" {
			t.Skip = true
			return nil
		}
		t.Content = os.Getenv(t.Name)
		if t.Content == "" {
			t.Content = t.Default
		}
		return next.UnmarshalField(ctx, field)
	})
}
