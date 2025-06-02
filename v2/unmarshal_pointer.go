package envy

import "reflect"

type _pointer reflect.Value

func (p _pointer) UnmarshalText(text []byte) (err error) {
	value := reflect.Value(p)
	temp := reflect.New(value.Type().Elem())
	err = _struct(reflect.Indirect(temp)).UnmarshalText(text)
	if err != nil {
		return err
	}
	if !reflect.Indirect(temp).IsZero() {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value.Set(temp)

	}

	return
}
