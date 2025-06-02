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
		match_expression := field.Tag.Get(matches_tagname)
		if match_expression == "" {
			return next.UnmarshalField(ctx, field)
		}
		if t.Matcher, err = regexp.Compile(match_expression); err != nil {
			return err
		}
		if t.Matcher.Match([]byte(t.Content)) {
			return next.UnmarshalField(ctx, field)
		}
		return DoesNotMatchError(field.Name, match_expression, t.Content)
	})
}
