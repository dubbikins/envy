package envy

import (
	"encoding"
	"reflect"
)

type FieldReflection struct {
	reflect.StructField
	Value ValueReflection
}

func NewFieldReflection(f reflect.StructField, v ValueReflection) (fr *FieldReflection, err error) {
	fr = &FieldReflection{f, v}
	return
}

func (field *FieldReflection) Indirect() {
	newValue := reflect.Indirect(field.Value.Value)
	field.Value = ValueReflection{reflect.Indirect(field.Value.Value), newValue.Kind()}
}

func (f *FieldReflection) Ref() any {
	return f.Value.Addr().Interface()
}

func (f *FieldReflection) Set(ref any) {
	f.Value.Set(reflect.Indirect(reflect.ValueOf(ref)))
}

func (field *FieldReflection) Unmarshal() (err error) {
	tag, err := field.Tag()
	if err != nil {
		return
	}
	if !field.Value.IsValid() {
		return INVALID_FIELD_ERROR
	}
	ref := field.Value.Addr().Interface()
	if unmarshaller, ok := ref.(encoding.TextUnmarshaler); ok {
		return unmarshaller.UnmarshalText([]byte(tag.Value))
	}
	if unmarshaller, ok := ref.(encoding.BinaryUnmarshaler); ok {
		return unmarshaller.UnmarshalBinary([]byte(tag.Value))
	}
	switch field.Value.Kind {
	case reflect.Ptr:
		if field.Value.IsNil() {
			field.Value.Set(reflect.New(field.Value.Type().Elem()))
		}
		if field.Value.Type().Elem().Kind() != reflect.Struct {
			field.Indirect()
			return field.Unmarshal()
		}
		err = Unmarshal(field.Value.Interface())
	case reflect.Struct:

		if err = Unmarshal(ref); err != nil {
			return err
		}
		field.Value.Set(reflect.Indirect(reflect.ValueOf(ref)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		UnmarshalInt(tag, field.Value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		UnmarshalUint(tag, field.Value)
	case reflect.String:
		UnmarshalString(tag, field.Value)
	case reflect.Bool:
		UnmarshalBool(tag, field.Value)
	case reflect.Float32, reflect.Float64:
		UnmarshalFloat(tag, field.Value)
	}
	return
}
