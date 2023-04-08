package envy

import (
	"log"
	"reflect"
)

type FieldReflection struct {
	reflect.StructField
	Value  ValueReflection
	Reader EnvironmentReader
}

func NewFieldReflection(f reflect.StructField, v ValueReflection, reader EnvironmentReader) (fr *FieldReflection, err error) {
	log.Printf("Creating Field Reflection: %s\n", f.Name)

	fr = &FieldReflection{f, v, reader}
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
	log.Printf("Unmarshalling Field: %s\n", field.Name)
	log.Printf("Field Kind: %s\n", field.Type.Kind())
	log.Printf("Field Name: %s\n", field.Name)
	log.Printf("Field Value: %s\n", field.Value)
	log.Printf("Field Interface: %s\n", field.Value.Interface())
	if !field.Value.IsValid() {
		return INVALID_FIELD_ERROR
	}
	switch field.Value.Kind {
	case reflect.Ptr:

		log.Printf("Field Kind: %s\n", field.Value.Type().Elem().Kind())
		if field.Value.IsNil() {
			field.Value.Set(reflect.New(field.Value.Type().Elem()))
		}
		if field.Value.Type().Elem().Kind() != reflect.Struct {
			field.Indirect()
			return field.Unmarshal()
		}
		err = Unmarshal(field.Value.Interface(), field.Reader)
	case reflect.Struct:
		ref := field.Value.Addr().Interface()
		if unmarshaller, ok := ref.(Unmarshaller); ok {
			return unmarshaller.Unmarshal(field)
		}
		if err = Unmarshal(ref, field.Reader); err != nil {
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
