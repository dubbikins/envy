package envy

import (
	"context"
	"errors"
	"reflect"
)

func Unmarshal(s any, options ...func(*Options)) (err error) {
	if reflect.TypeOf(s).Kind() != reflect.Pointer {
		return errors.New("unmarshalling reflection error: value passed to Unmarshal must be a struct pointer type")
	}
	opts := default_options
	for _, option := range options {
		option(opts)
	}
	ctx := WithOptionsContext(context.Background(), opts)
	element := reflect.TypeOf(s).Elem()
	parent_value := reflect.ValueOf(s)
	for i := 0; i < element.NumField(); i++ {
		// var tag *tag
		field := element.Field(i)
		if field.IsExported() {
			value := reflect.ValueOf(s).Elem().Field(i)
			tag, err := NewTag(value, parent_value)
			if err = tag.UnmarshalField(ctx, field); err != nil {
				return err
			}
		}
	}
	return
}

// (f, options ...func(tag TagMiddleware)) (err error) {
// 	if reflect.TypeOf(s).Kind() != reflect.Pointer {
// 		return errors.New("unmarshalling reflection error: value passed to Unmarshal must be a struct pointer type")
// 	}
// 	element := reflect.TypeOf(s).Elem()
// 	for i := 0; i < element.NumField(); i++ {
// 		// var tag *tag
// 		field := element.Field(i)
// 		if field.IsExported() {
// 			value := reflect.ValueOf(s).Elem().Field(i)
// 			tag := NewTag(value)
// 			for _, option := range options {
// 				option(tag)
// 			}
// 			if err = tag.UnmarshalField(context.Background(), field); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return
// }
