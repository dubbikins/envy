package envy

import (
	"context"
	"html/template"
	"io"
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
		if required_tag = field.Tag.Get(required_tagname); required_tag == "" {
			return next.UnmarshalField(ctx, field)
		}
		tmpl := template.Must(template.New("required").Parse(required_tag))
		tmpl.Execute(t, t.Parent.Interface())
		parsed, err := io.ReadAll(t)
		if t.Required, err = strconv.ParseBool(string(parsed)); err != nil {
			return err
		}
		if t.Content == "" && t.Required {
			return RequiredError(t.Name, t.Contents())
		}
		return next.UnmarshalField(ctx, field)
	})
}
