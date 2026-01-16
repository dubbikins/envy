package _default

import (
	"log/slog"
	"os"
	"text/template"

	"github.com/dubbikins/envy/v2/tag"
)

var TemplateFunctions = template.FuncMap{
	"env": func(name string) string {
		return os.Getenv(name)
	},
}


func WalkFn(next tag.WalkFn) tag.WalkFn {
	return func(node *tag.Node) (err error) {
		if !node.Value().IsZero() {
			return next(node)
		}
		if node.Field() == nil { 
			return 
		}
		if err = node.Parse("default", tag.LexTemplate); err != nil || node.Skipped(){
			return 
		}
		for _, value := range node.TagValues() {
			if value != "" {
				
				var _tmpl = template.New("default").Funcs(TemplateFunctions)
				if _tmpl, err = _tmpl.Parse(value); err != nil {
					return
				}
				if err = _tmpl.Execute(node, node.Parent().Value().Interface()); err != nil {
					slog.Error("default tag template error")
					return
				}
				if err = node.UnmarshalText([]byte(node.Bytes())); err != nil {
					return
				}
				return next(node)
			}
		}
		return next(node)
	}
}
