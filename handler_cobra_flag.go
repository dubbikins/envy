package envy

// import (
// 	"context"
// 	"log/slog"
// 	"reflect"
// 	"strings"

// 	"github.com/spf13/cobra"
// )

// const flag_tagname = "flag"

// func WithCobraFlag(cmd *cobra.Command) func (next TagHandler) TagHandler {
// 	return func(next TagHandler) TagHandler {
// 		return TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {

// 		t, err := GetTagContext(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		t.Raw = field.Tag.Get(flag_tagname)
// 		segments := strings.Split(t.Raw, ";")
// 		if len(segments) <= 0 {
// 			slog.Error("unmarshaling cobra flag", "raw", t.Raw)
// 			return next.UnmarshalField(ctx, field)
// 		}

// 		t.Name = segments[0]

// 		if t.Skip {
// 			slog.Debug("env tag found, but value is already set when OverrideValues is false, skipping", "tag", field.Tag.Get(flag_tagname))
// 			return next.UnmarshalField(ctx, field)
// 		}else if t.Name == "-" {
// 			slog.Debug("env tag found, but value is already set when OverrideValues is false, skipping", "tag", field.Tag.Get(flag_tagname))

// 			t.Skip = true
// 			return next.UnmarshalField(ctx, field)
// 		}
// 		var flag = cmd.Flag(t.Name)
// 		//slog.Info("unmarshaling cobra flag", "flag",flag, "name", t.Name)
// 		// if flag == nil {
// 		// 	slog.Error("flag not found", "name", t.Name)
// 		// 	return fmt.Errorf("flag '%s' not found by cmd.Flag", t.Name)
// 		// }
// 		t.Content = flag.Value.String()
// 		slog.Debug("Unmarshalling from cobra.Command.Flag", "name", t.Name, "value", t.Content )
// 		return next.UnmarshalField(ctx, field)
// 	})
// 	}
// }