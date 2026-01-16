package _default

import (
	"testing"

	"github.com/dubbikins/envy/v2"
)


func TestDefaultWalkFn(t *testing.T) {
	type Example struct {
		Field string `default:"test"`
		TemplateField string `default:"{{.Field}}-template-{{env \"TEST\"}}"`
	}
	t.Setenv("TEST", "ok")
	var err error
	have := &Example{}
	expected := Example{
		Field: "test",
		TemplateField: "test-template-ok",
	}
	if err = envy.Unmarshal(have, WalkFn); err != nil {
		t.Fatal(err)
	}

	if have.Field != expected.Field {
		t.Fatalf("Expected .Field to be %q but was %q", expected.Field, have.Field)
	}
	if have.TemplateField != expected.TemplateField {
		t.Fatalf("Expected .TemplateField to be %q but was %q", expected.TemplateField, have.TemplateField)
	}
}