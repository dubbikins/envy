package envy

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

type _int reflect.Value

func (i _int) UnmarshalText(text []byte) (err error) {
	var value = reflect.Value(i)
	var bitBase int
	switch value.Kind() {
	case reflect.Int:
		bitBase = 0
	case reflect.Int8:
		bitBase = 8
	case reflect.Int16:
		bitBase = 16
	case reflect.Int32:
		bitBase = 32
	default:
		bitBase = 64
	}
	if val, err := strconv.ParseInt(strings.ReplaceAll(string(text), ",", ""), 0, bitBase); err != nil {
		log.Printf("error unmarshalling int[%s]: %s\n", string(text), err.Error())
		return err
	} else {
		value.SetInt(val)
	}
	return
}
