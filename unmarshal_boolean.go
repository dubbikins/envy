package envy

import (
	"reflect"
	"strconv"
)

type _boolean reflect.Value

func (b _boolean) UnmarshalText(text []byte) (err error) {
	var val bool
	switch string(text) {
	case "yes", "on":
		val = true
	case "no", "off":
		val = false
	default:
		if val, err = strconv.ParseBool(string(text)); err != nil {
			return
		}
	}
	reflect.Value(b).SetBool(val)
	return
}
