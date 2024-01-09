package envy_test

import (
	"fmt"
	"html/template"
	"reflect"
	"testing"
)

type Writer struct {
	value []byte
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.value = p
	fmt.Println(string(p))
	return len(p), nil
}
func TestXxx(t *testing.T) {
	type Test struct {
		Name string `default:"not_test"`
	}

	value := &Test{
		Name: "test",
	}
	rv := reflect.ValueOf(value)

	tmpl := template.Must(template.New("test").Parse(`{{eq .Name "test"}}`))
	w := &Writer{}
	tmpl.Execute(w, rv.Interface())
	expected := "true"
	have := string(w.value)
	if have != expected {
		t.Fatalf("Expected value to be %q but was %q", expected, have)
	}
}
