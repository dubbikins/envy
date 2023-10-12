package envy

import (
	"context"
	"reflect"
	"strconv"
)

const required_tagname = "required"

func WithRequiredTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {

		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		var required_tag string
		if required_tag = string(field.Tag.Get(required_tagname)); required_tag == "" {
			return next.UnmarshalField(ctx, field)
		}
		t.Required, err = strconv.ParseBool(required_tag)
		if t.Value == "" && t.Required {
			return RequiredError(t.Name, t.Contents())
		}
		return next.UnmarshalField(ctx, field)
	})
}
