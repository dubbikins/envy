package envy

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

const env_tagname = "env"


func WithEnvTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		t.Name = field.Tag.Get(env_tagname)
		//If handling the default tag set skip, we want to still apply the remaining handlers like required, options, etc
		// so we return next.UnmarshalField(ctx, field), but if the name is "-", we want to skip the rest of the handlers
		// so we return nil
		if t.Skip {
			slog.Debug("env tag found, but value is already set when OverrideValues is false, skipping", "tag", field.Tag.Get(env_tagname))
			return next.UnmarshalField(ctx, field)
		}
		if t.Name == "-" {
			t.Skip = true
			return nil
		}
		for _, tag := range strings.Split(t.Name, ";") {
			t.Content = os.Getenv(tag)
			if t.Content != "" {
				t.SelectedName = tag
				break
			}
		}
		
		if t.Content == "" {
			t.Content = t.Default
		}
		return next.UnmarshalField(ctx, field)
	})
}
