package envy

import (
	"context"
	"reflect"
	"strings"
)

const options_tagname = "options"

func WithOptionsTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {

		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		var rawTag string
		if rawTag = string(field.Tag.Get(options_tagname)); rawTag == "" {
			return next.UnmarshalField(ctx, field)
		}
		t.Raw += rawTag
		if t.Options = strings.Split(strings.Trim(rawTag, "[({})]"), ","); len(t.Options) < 1 {
			return next.UnmarshalField(ctx, field)
		}
		for _, option := range t.Options {
			if t.Content == option {
				return next.UnmarshalField(ctx, field)
			}
		}
		return InvalidOptionError(t.Content, t.Options)
	})
}
