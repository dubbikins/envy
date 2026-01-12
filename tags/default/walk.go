package default_tag

import (
	"log/slog"
	"os"
	"text/template"

	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/tag/text"
)

var TemplateFunctions = template.FuncMap{
	"env": func(name string) string {
		return os.Getenv(name)
	},
}


func WalkFn(next tag.WalkFn) tag.WalkFn {
	return func(node *tag.Node) (err error) {
		if node.Field() == nil { //&& node.Value().Kind() != reflect.Pointer
			return 
		}
		if err = text.Parse("default", node, text.LexDefaultWithTemplate); err != nil || node.Skipped(){
			return 
		}
		for _, value := range node.TagValues() {
			if value != "" {
				var _tmpl = template.New("default").Funcs(TemplateFunctions)
				slog.Info("executing template", "template", value, "data", node.Parent().Value().Interface(),"tag", node.Field().Tag )

				if _tmpl, err = _tmpl.Parse(value); err != nil {
					return
				}
				if err = _tmpl.Execute(node, node.Parent().Value().Interface()); err != nil {
					slog.Error("default tag template error")
					return
				}
				slog.Info("Node Write", "bytes", node.Bytes())
				if err = node.UnmarshalText([]byte(node.Bytes())); err != nil {
					return
				}
				slog.Info("Node Unmarshalled", "bytes", node.Value())

				return next(node)
			}else {
				slog.Info("Empty tag value", "field", node.Field().Name)
			}
		}
		slog.Info("no default template", "field", node.Field().Name, "tag", node.Field().Tag	)

		return next(node)
	}
}
