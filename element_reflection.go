package envy

import (
	"reflect"
)

type ElementReflection struct {
	reflect.Type
	Fields []*FieldReflection
}

func NewElementReflection(t TypeReflection, v ValueReflection) (e *ElementReflection, err error) {
	e = &ElementReflection{t.Elem(), []*FieldReflection{}}

	for i := 0; i < e.NumField(); i++ {
		var fr *FieldReflection
		if fr, err = NewFieldReflection(e.Field(i), ValueReflection{v.Elem().Field(i), v.Elem().Field(i).Kind()}); err == nil {

			e.Fields = append(e.Fields, fr)
		} else {
			return e, err
		}
	}
	return
}
