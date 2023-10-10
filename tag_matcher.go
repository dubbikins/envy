package envy

import (
	"context"
	"reflect"
	"regexp"
)

const matches_tagname = "matches"

func WithMatchesTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {

		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		rawTag := field.Tag.Get(matches_tagname)
		if rawTag == "" {
			return next.UnmarshalField(ctx, field)
		}
		if t.Matcher, err = regexp.Compile(rawTag); err != nil {
			return err
		}
		if !t.Matcher.Match([]byte(t.Value)) {
			return DoesNotMatchError(t.Name, t.Value, rawTag)
		}
		return InvalidOptionError(t.Value, t.Options)
	})
}
