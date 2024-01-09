package envy

import (
	"context"
	"html/template"
	"reflect"
	"strconv"
)

type Writer struct {
	value []byte
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.value = p
	return len(p), nil
}

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
		w := &Writer{}
		tmpl.Execute(w, t.Parent.Interface())
		t.Required, err = strconv.ParseBool(required_tag)

		if t.Content == "" && t.Required {
			return RequiredError(t.Name, t.Contents())
		}
		return next.UnmarshalField(ctx, field)
	})
}
