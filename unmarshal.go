package envy

import (
	"log"
)

type Unmarshaller interface {
	Unmarshal(*FieldReflection) error
	// Unmarshals() (reflect.Kind, string)
	// Indirect()
}

func Unmarshal(s any, reader EnvironmentReader) (err error) {
	var r *Reflection
	r, err = NewReflection(s, reader)
	if err != nil {
		return
	}
	log.Printf("Field Reflections: %v \n", r.Element.Fields)
	for _, field := range r.Element.Fields {
		err = field.Unmarshal()
		if err != nil {
			return
		}
	}
	return err

}

type UnmarshallerOptions struct {
	EnvironmentReader EnvironmentReader
}

type DefaultPrimitiveUnmarshaller struct {
}

func (u *DefaultPrimitiveUnmarshaller) Unmarshal(...EnvironmentReader) (err error) {

	return
}
