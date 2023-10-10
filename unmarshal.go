package envy

import (
	"encoding"
	"errors"
	"reflect"
)

func FromEnvironmentAs[T any](t *T) error { return Unmarshal(t) }

func Unmarshal(s any) (err error) {
	if reflect.TypeOf(s).Kind() != reflect.Pointer {
		return errors.New("unmarshalling reflection error: value passed to Unmarshal must be a struct pointer type")
	}
	element := reflect.TypeOf(s).Elem()
	for i := 0; i < element.NumField(); i++ {
		var ok bool
		var tag *tag
		var unmarshaller encoding.TextUnmarshaler
		field := element.Field(i)
		if field.IsExported() {
			value := reflect.ValueOf(s).Elem().Field(i)
			if tag, err = getTag(field); err != nil {
				return err
			} else if !value.IsValid() {
				return INVALID_FIELD_ERROR
			}
			ref := value.Addr().Interface()
			unmarshaller, ok = ref.(encoding.TextUnmarshaler)
			if !ok {
				switch value.Kind() {
				case reflect.Ptr:
					unmarshaller = _pointer(value)
				case reflect.Struct:
					unmarshaller = _struct(value)
				case reflect.Slice:
					unmarshaller = _slice(value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					unmarshaller = _int(value)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					unmarshaller = _uint(value)
				case reflect.String:
					unmarshaller = _string(value)
				case reflect.Bool:
					unmarshaller = boolean(value)
				case reflect.Float32, reflect.Float64:
					unmarshaller = _float(value)
				}
			}
			err = unmarshaller.UnmarshalText([]byte(tag.Value))
			if err != nil {
				return err
			}
		}
	}
	return

}
