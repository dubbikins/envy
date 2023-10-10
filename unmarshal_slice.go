package envy

import (
	"fmt"
	"reflect"
)

type _slice reflect.Value

func (s _slice) UnmarshalText(text []byte) (err error) {

	return fmt.Errorf("slice unmarshal not implemented")
}
