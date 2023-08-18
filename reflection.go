package envy

import (
	"errors"
	"reflect"
)

type Reflection struct {
	Value   ValueReflection
	Type    TypeReflection
	Element *ElementReflection
}
type ValueReflection struct {
	reflect.Value
	Kind reflect.Kind
}

type TypeReflection struct {
	reflect.Type
	Kind reflect.Kind
}

func NewReflection(s interface{}) (r *Reflection, err error) {
	r = &Reflection{
		Type:  TypeReflection{reflect.TypeOf(s), reflect.TypeOf(s).Kind()},
		Value: ValueReflection{reflect.ValueOf(s), reflect.ValueOf(s).Kind()},
	}
	if r.Type.Kind != reflect.Pointer {
		return nil, errors.New("unmarshalling reflection error: value passed to Unmarshal must be a struct pointer type")
	}
	r.Element, err = NewElementReflection(r.Type, r.Value)
	return r, err
}
