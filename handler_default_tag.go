package envy

import (
	"context"
	"html/template"
	"io"
	"reflect"
)

const default_tagname = "default"

func WithDefaultTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		var tag_value = field.Tag.Get(default_tagname)
		tmpl := template.Must(template.New("default").Parse(tag_value))
		tmpl.Execute(t, t.Parent.Interface())
		parsed, err := io.ReadAll(t)
		t.Default = string(parsed)
		return next.UnmarshalField(ctx, field)
	})
}
