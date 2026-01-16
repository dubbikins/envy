package dynamodb

import (
	"log/slog"

	"github.com/dubbikins/envy/v2/tag"
)


func WalkFn( next tag.WalkFn) tag.WalkFn {
	return func( node *tag.Node) (err error) {
		
		//If this isn't a struct field element, then we should stop walking the tag
		if node.Field() == nil { 
			return 
		}
		if err = node.Parse("dynamodbav", lexTag); err != nil || node.Skipped(){
			return 
		}
	
		//setup 		
			
		for _, key := range node.TagValues() {
			slog.Info(key)
		}
		return next(node)
	}
}
