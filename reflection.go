package envy

import (
	"errors"
	"log"
	"reflect"
)

type Reflection struct {
	Value   ValueReflection
	Type    TypeReflection
	Element *ElementReflection
	Reader  EnvironmentReader
}
type ValueReflection struct {
	reflect.Value
	Kind reflect.Kind
}

type TypeReflection struct {
	reflect.Type
	Kind reflect.Kind
}

func NewReflection(s interface{}, reader EnvironmentReader) (r *Reflection, err error) {
	r = &Reflection{
		Type:   TypeReflection{reflect.TypeOf(s), reflect.TypeOf(s).Kind()},
		Value:  ValueReflection{reflect.ValueOf(s), reflect.ValueOf(s).Kind()},
		Reader: reader,
	}
	if r.Type.Kind != reflect.Pointer {
		return nil, errors.New("unmarshalling reflection error: value passed to Unmarshal must be a struct pointer type")
	}
	log.Println("Unmarshalling Interface...")
	log.Println("Type: " + r.Type.String())
	log.Println("Kind: " + r.Type.Kind.String())
	log.Println("Value: " + r.Value.String())
	// log.Printf("Can Address Value: %t\n", r.Value.CanAddr())
	// log.Printf("Can Set Value: %t\n", r.Value.CanSet())
	log.Println("Unmarshalling Element...")
	r.Element, err = NewElementReflection(r.Type, r.Value, reader)
	return r, err
}
