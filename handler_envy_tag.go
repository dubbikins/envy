package envy

import (
	"context"
	"log/slog"
	"reflect"
	"strings"
)

const envy_global_tagname = "envy"

func WithEnvyGlobalTag(next TagHandler) TagHandler {
	return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) (err error) {
		var tag_value = field.Tag.Get(envy_global_tagname)
		
		t, err := GetTagContext(ctx)
		if err != nil {
			return err
		}
		if t.tag_unmarshaller_opts == nil {
			t.tag_unmarshaller_opts = &tagUnmarshallerOptions{}
		}
		if tag_value == "" {
			slog.Debug("no tag value found, skipping")
			return next.UnmarshalField(ctx, field)
		}
		slog.Debug("setting up envy options", "tag", tag_value)
		tag_segments := strings.Split(field.Tag.Get(envy_global_tagname), ";")
		for _, segment := range tag_segments {
			if strings.HasPrefix(segment, "@") {
				if err = t.tag_unmarshaller_opts.Load(segment); err != nil {
					return err
				}
			} else {
				slog.Debug("unknown option when setting up envy options", "tag", segment)
			}

		}
		
		return next.UnmarshalField(ctx, field)
	})
}
