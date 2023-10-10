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
			//No options were defined, so we don't need to do validations against them
			//We can just return the next unmarshaller
			return next.UnmarshalField(ctx, field)
		}
		for _, option := range t.Options {
			if t.Value == option {
				return next.UnmarshalField(ctx, field)
			}
		}
		return InvalidOptionError(t.Value, t.Options)
	})
}
