package envy

import (
	"fmt"
	"reflect"
)

type _slice reflect.Value

func (s _slice) UnmarshalText(text []byte) (err error) {
	// value := reflect.Value(s)
	// for i := 0; i < value.Len(); i++ {
	// 	if err = Unmarshal(value.Index(i).Addr().Interface()); err != nil {
	// 		return err
	// 	}
	// }
	return fmt.Errorf("slice unmarshal not implemented")
}
