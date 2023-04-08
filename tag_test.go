package envy

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"testing"
)

type Test struct {
	String      string `env:"string"`
	Int         int    `env:"int"`
	Int8        int    `env:"int8"`
	Int16       int    `env:"int16"`
	Int32       int    `env:"int32"`
	Int64       int    `env:"int64"`
	Uint        int    `env:"uint"`
	Uint8       int    `env:"uint8"`
	Uint16      int    `env:"uint16"`
	Uint32      int    `env:"uint32"`
	Uint64      int    `env:"uint64"`
	Bool        bool   `env:"bool"`
	BoolDefault bool   `env:"bool"`
	Pointer     *struct {
		Field string `env:"pointer_field"`
	}
	Struct struct {
		Field string `env:"struct_field"`
	}
	MultiFieldStruct struct {
		Field1 string `env:"multi_field_struct_1"`
		Field2 string `env:"multi_field_struct_2"`
	}
}

type TestError struct {
	Bool struct {
		Default bool `env:"bool;default=true"`
	}
}

type TC struct {
	expected    any
	envar_key   string
	envar_value string
	field       string
	format      string
	uut         *Test
	value       func(uut *Test) any
}

func (tc *TC) SetupTest(t *testing.T) func(tb *testing.T) {
	log.Println("setup test case")
	t.Setenv(tc.envar_key, tc.envar_value)
	//tb.Setenv("", "")
	tc.uut = &Test{}
	err := Unmarshal(tc.uut, &OSEnvironmentReader{})
	if err != nil {
		t.Fatal(err)
	}
	// Return a function to teardown the test
	return func(tb *testing.T) {
		log.Println("teardown suite")
	}
}

func (tc *TC) Name() string {
	return fmt.Sprintf("Test field %s set from env var %s with %s.", tc.field, tc.envar_key, tc.envar_value)
}
func (tc *TC) CheckError(t *testing.T) {
	if tc.value(tc.uut) != tc.expected {
		t.Fatalf("expected %s field to equal '"+tc.format+"', but was '"+tc.format+"'", tc.field, tc.expected, tc.value)
	}
}
func (tc *TC) Run(t *testing.T) {

	t.Run(tc.Name(), func(t *testing.T) {
		teardownTest := tc.SetupTest(t)
		defer teardownTest(t)
		tc.CheckError(t)
	})
}

var test_cases = []*TC{
	{
		field:       "String",
		envar_key:   "string",
		envar_value: "test",
		expected:    "test",
		format:      "%s",
		value:       func(uut *Test) any { return uut.String },
	},
	{
		field:       "Int",
		envar_key:   "int",
		envar_value: "-2147483648",
		expected:    -2147483648,
		format:      "%d",
		value:       func(uut *Test) any { return uut.Int },
	},
	{
		field:       "Int",
		envar_key:   "int",
		envar_value: "2,147,483,647", //can parse numeric values with commas
		expected:    2147483647,
		format:      "%d",
		value:       func(uut *Test) any { return uut.Int },
	},

	{
		field:       "Int8",
		envar_key:   "int8",
		envar_value: "127",
		expected:    127,
		format:      "%d",
		value:       func(uut *Test) any { return uut.Int8 },
	},
	{
		field:       "Int16",
		envar_key:   "int16",
		envar_value: "127",
		expected:    127,
		format:      "%d",
		value:       func(uut *Test) any { return uut.Int16 },
	},
	{
		field:       "Int16",
		envar_key:   "int16",
		envar_value: "127",
		expected:    127,
		format:      "%d",
		value:       func(uut *Test) any { return uut.Int16 },
	},
	{
		field:       "Bool",
		envar_key:   "bool",
		envar_value: "true",
		expected:    true,
		format:      "%b",
		value:       func(uut *Test) any { return uut.Bool },
	},
	{
		field:       "Struct.Field",
		envar_key:   "struct_field",
		envar_value: "test",
		expected:    "test",
		format:      "%s",
		value:       func(uut *Test) any { return uut.Struct.Field },
	},
	{
		field:       "Pointer.Field",
		envar_key:   "pointer_field",
		envar_value: "test",
		expected:    "test",
		format:      "%s",
		value:       func(uut *Test) any { return uut.Pointer.Field },
	},
	{
		field:       "MultiFieldStruct.Field1",
		envar_key:   "multi_field_struct_1",
		envar_value: "test1",
		expected:    "test1",
		format:      "%s",
		value:       func(uut *Test) any { return uut.MultiFieldStruct.Field1 },
	},
	{
		field:       "MultiFieldStruct.Field2",
		envar_key:   "multi_field_struct_2",
		envar_value: "test2",
		expected:    "test2",
		format:      "%s",
		value:       func(uut *Test) any { return uut.MultiFieldStruct.Field2 },
	},
}

func TestAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	log.SetOutput(w)
	for _, tc := range test_cases {
		tc.Run(t)
	}
}

type Test2 struct {

	// Pointer     *struct {
	// 	Field string `env:"pointer_field"`
	// }
	Struct struct {
		Field *string `env:"struct_field"`
	}
	// MultiFieldStruct struct {
	// 	Field1 string `env:"string"`
	// 	Field2 string `env:"string"`
	// }
}
