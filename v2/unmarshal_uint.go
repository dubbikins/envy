package envy

import (
	"reflect"
	"strconv"
	"strings"
)

type _uint reflect.Value

func (u _uint) UnmarshalText(text []byte) (err error) {
	var value = reflect.Value(u)
	var bitsize int
	switch value.Kind() {
	case reflect.Uint:
		bitsize = 0
	case reflect.Uint8:
		bitsize = 8
	case reflect.Uint16:
		bitsize = 16
	case reflect.Uint32:
		bitsize = 32
	default:
		bitsize = 64
	}
	if val, err := strconv.ParseUint(strings.ReplaceAll(string(text), ",", ""), 0, bitsize); err != nil {
		return err
	} else {
		value.SetUint(val)
	}
	return
}
