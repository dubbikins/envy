package default_tag

import (
	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/tag/text"
)


func WalkFn(next tag.WalkFn) tag.WalkFn {
	return func(node *tag.Node) (err error) {
		if node.Field() == nil { //&& node.Value().Kind() != reflect.Pointer
			return 
		}
		if err = text.Parse("default", node, text.LexEnvironmentVariableTag); err != nil || node.Skipped(){
			return 
		}
		for _, value := range node.TagValues() {
			if value != "" {
				if err = node.UnmarshalText([]byte(value)); err != nil {
					return
				}
				break
			}
		}
		return next(node)
	}
}
