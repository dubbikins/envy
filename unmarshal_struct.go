package envy

import (
	"reflect"
)

type _struct reflect.Value

func (s _struct) UnmarshalText(text []byte) (err error) {
	value := reflect.Value(s)
	ref := value.Addr().Interface()
	if err = Unmarshal(ref); err != nil {
		return err
	}
	value.Set(reflect.Indirect(reflect.ValueOf(ref)))
	return
}
