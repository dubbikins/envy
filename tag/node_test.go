package tag

import (
	"reflect"
	"testing"
)


func TestSliceElementValues(t *testing.T) {
	type Example struct {
		STRING_SLICE []*string `env:"STRING_SLICE"`
	}
	example := &Example{}
	n := NewRootNode( example)
	n = n.Descendant(n.value.Elem(), nil) //pointer -> struct

	var element = n.value.Type()
	var field = element.Field(0)
	var value = reflect.Value(n.value).Field(0)
	n = n.Descendant(value,&field)
	n.separator = byte(',')
	n.content.WriteString("1,2,3")
	if n.value.Kind() != reflect.Slice{
		t.Fatalf("expected slice but was %s", n.value.Kind())
	}
	var want_elements = []string{"1", "2", "3"}
	var i int
	for elem := range n.slice_element_values([]byte("1,2,3")) {
		if string(elem) != want_elements[i] {
			t.Fatalf("expected element[%d] to be '%s' but was '%s'", i, want_elements[i], string(elem))
		}
		i++
	}

}