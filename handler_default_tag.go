package envy

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"reflect"
)

const default_tagname = "default"

func WithDefaultTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		if t == nil {
			return fmt.Errorf("tag context is nil")
		}
		if t.tag_unmarshaller_opts == nil {
			t.tag_unmarshaller_opts = &tagUnmarshallerOptions{}
		}
		if !t.Value.IsZero() && !t.tag_unmarshaller_opts.OverrideValues {
			slog.Debug("default tag found, but value is already set when OverrideValues is false, skipping", "tag", field.Tag.Get(default_tagname))
			t.Skip = true
			return next.UnmarshalField(ctx, field)
		}
		var tag_value = field.Tag.Get(default_tagname)
		tmpl, err := template.New("default").Parse(tag_value)
		if err != nil {
			return err
		}
		if err = tmpl.Execute(t, t.Parent.Interface()); err != nil {
			return err
		}
		parsed, err := io.ReadAll(t)
		if err != nil {
			return err 
		}
		t.Default = string(parsed)
		return next.UnmarshalField(ctx, field)
	})
}
