package envy

import (
	"reflect"
	"strconv"
	"strings"
)

type _float reflect.Value

func (f _float) UnmarshalText(text []byte) (err error) {
	var value = reflect.Value(f)
	bitsize := 0
	switch value.Kind() {
	case reflect.Float32:
		bitsize = 32
	case reflect.Float64:
		bitsize = 64
	}
	if val, err := strconv.ParseFloat(strings.ReplaceAll(string(text), ",", ""), bitsize); err == nil {
		value.SetFloat(val)
	}
	return
}
