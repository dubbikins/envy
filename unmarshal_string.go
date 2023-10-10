package envy

import "reflect"

type _string reflect.Value

func (s _string) UnmarshalText(text []byte) (err error) {
	reflect.Value(s).SetString(string(text))
	return
}
