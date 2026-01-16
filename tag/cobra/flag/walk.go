package flag

import (
	"fmt"

	"github.com/dubbikins/envy/v2/tag"
	"github.com/spf13/cobra"
)


type Node struct {
	*tag.Node
}

func (n *Node) Set(value string) (err error) {
	return n.Node.UnmarshalText([]byte(value))
}

func (n *Node)Type() (_type string){
	var ok bool
	if _type, ok = n.Option("type"); ok {
		return
	}
	return n.Value().Type().String()
}
const TAGNAME = "flag"

func TagWalkFn(cmd *cobra.Command) func(next tag.WalkFn) tag.WalkFn {
	return func(next tag.WalkFn) tag.WalkFn {
		return func(node *tag.Node) (err error) {
		//If this isn't a struct field element, then we should stop walking the tag
		if node.Field() == nil {
			return 
		}
		if err = node.Parse(TAGNAME, LexTag); err != nil {
			return
		}
		return next(node)
		}
	}
}

func InitCmd(cmd *cobra.Command) tag.WalkFn {
	return func(node *tag.Node) (err error) {
		//If this isn't a struct field element, then we should stop walking the tag
		if node.Field() == nil {
			return 
		}
		
		if err = node.Parse(TAGNAME, LexTag); err != nil || node.Skipped() {
			return
		}
		var name, shorthand, usage string

		for _, v := range node.TagValues() {
			if name == "" {
				name = v
			}else if shorthand == "" {
				shorthand = v
			}else {
				err = fmt.Errorf("invalid 'flag' tag syntax: cannot have more than name and shorthand")
				return
			}
		}
		var found bool
		if usage, found = node.Option("usage"); !found {
			usage, _ = node.Option("description")
		}
		cmd.PersistentFlags().VarP(&Node{node},name, shorthand, usage)
		return
		}
}

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

