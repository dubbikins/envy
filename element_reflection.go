package envy

import (
	"log"
	"reflect"
)

type ElementReflection struct {
	reflect.Type
	Fields []*FieldReflection
}

func NewElementReflection(t TypeReflection, v ValueReflection) (e *ElementReflection, err error) {
	e = &ElementReflection{t.Elem(), []*FieldReflection{}}
	log.Println("Unmarshalling Element Field Values...")
	for i := 0; i < e.NumField(); i++ {
		var fr *FieldReflection
		if fr, err = NewFieldReflection(e.Field(i), ValueReflection{v.Elem().Field(i), v.Elem().Field(i).Kind()}); err == nil {
			log.Println("Appending Field Reflection")
			e.Fields = append(e.Fields, fr)
		} else {
			return e, err
		}
	}
	return
}
