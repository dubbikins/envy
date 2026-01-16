package envy

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dubbikins/envy/v2/types"
)

// #################################### TEST CASES ####################################
// #################################### Primitive Struct Fields ####################################
func TestUnmarshalStruct(t *testing.T) {
	//SETUP
	type Example struct {
		String string `env:"STR"`
		StringExpand string `env:"STRREF;expand"`
		StringExpandTemplate string `env:"STRREFSTRUCT;expand"`
		Int int `env:"INT"`
		Int8 int8 `env:"INT8"`
		Float32 float32 `env:"FLOAT32"`
		SkippedField string `env:"-"`
		Bool bool `env:"BOOL"`
		Complex64 complex64 `env:"COMPLEX64"`
		Complex128 complex128 `env:"COMPLEX128"`
		SkippedFieldWithEmptyFlags string `env:""`
		Required string `env:"REQUIRED" required:"true"`
		unexported string //should skip and be ok
	}

	var test_cases = StructTestCases[Example] {
		{
			name: "Comparable Type Fields",
			have: &Example{},
			want: &Example{
				String: "ok",
				StringExpand: "ok",
				StringExpandTemplate: "ok",
				Int: 10,
				Int8: 127,
				Float32: 12345.6789,
				Bool: true,
				Complex64: 1 + 2i,
				Required: "set",
			},
			setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
				t.Setenv("STR", tc.want.String)
				t.Setenv("STRREF", "$STR")
				t.Setenv("STRREFSTRUCT", "${.String}")
				t.Setenv("INT", fmt.Sprintf("%d", tc.want.Int))
				t.Setenv("INT8", fmt.Sprintf("%d", tc.want.Int8))
				t.Setenv("FLOAT32", fmt.Sprintf("%f", tc.want.Float32))
				t.Setenv("BOOL", fmt.Sprintf("%t", tc.want.Bool))
				t.Setenv("COMPLEX64", fmt.Sprintf("%v", tc.want.Complex64))
				t.Setenv("COMPLEX128", fmt.Sprintf("%v", tc.want.Complex128))
				t.Setenv("REQUIRED", "set")
			},
		},
		{
			name: "Boolean Yes",
			have: &Example{},
			want: &Example{
				Bool: true,
				Required: "set",
			},
			setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
				t.Setenv("BOOL", "yes")
				t.Setenv("REQUIRED", "set")
			},
		},
		{
			name: "Missing Required Field",
			have: &Example{},
			want: &Example{
			
			},
			expectedAnyError: true ,
			setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
		
				
			},
		},
		// {
		// 	name: "Boolean Off",
		// 	have: &Example{},
		// 	want: &Example{
		// 		Bool: false,
		// 	},
		// 	setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
		// 		t.Setenv("BOOL", "off")
		// 	},
		// },
		// {
		// 	name: "expands",
		// 	have: &Example{},
		// 	want: &Example{
		// 		String: "ok",
		// 		StringExpand: "ok",
		// 		StringExpandTemplate: "ok",
		// 		Int: 10,
		// 	},
		// 	setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
		// 		t.Setenv("STR", tc.want.String)
		// 		t.Setenv("INT", fmt.Sprintf("%d", tc.want.Int))
		// 	},
		// },
	}
	test_cases.Run(t)
}

// #################################### Slice Struct Fields ####################################
func TestUnmarshalStructWithStringSliceField(t *testing.T) {
	//SETUP
	
	type Example struct {
		SliceOfStrings []string `default:"test" env:"STR;sep=','"`
	}

	var test_cases = StructSliceTestCases[Example] {
		{
			name: "String Slice Example",
			have: &Example{},
			want: &Example{
				SliceOfStrings: []string{"foo", "bar"},
			},
			setEnv: func(t *testing.T, tc NonComparableTestCase[Example]) {
				t.Setenv("STR", strings.Join(tc.want.SliceOfStrings, ","))
				// t.Setenv("BYTES", "foobar")
			},
		},
	}
	test_cases.Run(t, func(t *testing.T, have, want Example) {
		if len(have.SliceOfStrings) != len(want.SliceOfStrings) {
			t.Fatalf("expected slice to have length '%d' but got '%d'", len(want.SliceOfStrings), len(have.SliceOfStrings))
		}
		for i := range want.SliceOfStrings {
			if want.SliceOfStrings[i] != have.SliceOfStrings[i] {
				t.Fatalf("expected slice String[%d] to be '%s' but got '%s'", i, want.SliceOfStrings[i], have.SliceOfStrings[i])
			}
		}
		
	})
}

func TestUnmarshalStructWithStringSliceFieldUsesSecondary(t *testing.T) {
	//SETUP
	
	type Example struct {
		Field string `env:"STR|STR2"`
	}

	var test_cases = StructTestCases[Example] {
		{
			name: "String Slice Example",
			have: &Example{},
			want: &Example{
				Field: "foo" ,
			},
			setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
				t.Setenv("STR", "foo")
				t.Setenv("STR2", "bar")
			},
		},
		{
			name: "String Slice Example",
			have: &Example{},
			want: &Example{
				Field: "bar" ,
			},
			setEnv: func(t *testing.T, tc ComparableTestCase[Example]) {
				//t.Setenv("STR", "foo")
				t.Setenv("STR2", "bar")
			},
		},
	}
	test_cases.Run(t)
}

func TestUnmarshalStructWithSliceFields(t *testing.T) {
	//SETUP
	type Subtype struct {
		String string `env:"BYTES"`
	}
	type Example struct {
		SliceOfStrings []string `env:"STR"`
		Bytes []byte `env:"BYTES;sep=''"`
		Ptr *Subtype
		Struct Subtype
		PtrIgnore *Subtype `env:"-"`
	}

	var test_cases = StructSliceTestCases[Example] {
		{
			name: "String Slice Example",
			have: &Example{},
			want: &Example{
				SliceOfStrings: []string{"foo", "bar"},
				Bytes: []byte("foobar"),
				Ptr: &Subtype{
					String: "foobar",
				},
				Struct: Subtype{
					String: "foobar",
				},
			},
			setEnv: func(t *testing.T, tc NonComparableTestCase[Example]) {
				t.Setenv("STR", strings.Join(tc.want.SliceOfStrings, ","))
				t.Setenv("BYTES", "foobar")
			},
		},
	}
	test_cases.Run(t, func(t *testing.T, have, want Example) {
		if len(have.SliceOfStrings) != len(want.SliceOfStrings) {
			t.Fatalf("expected slice to have length '%d' but got '%d'", len(want.SliceOfStrings), len(have.SliceOfStrings))
		}
		for i := range want.SliceOfStrings {
			if want.SliceOfStrings[i] != have.SliceOfStrings[i] {
				t.Fatalf("expected slice String[%d] to be '%s' but got '%s'", i, want.SliceOfStrings[i], have.SliceOfStrings[i])
			}
		}
		if len(have.Bytes) != len(want.Bytes) {
			t.Fatalf("expected byte slice to have length '%d' but got '%d'", len(want.Bytes), len(have.Bytes))
		}
		for i := range want.Bytes {
			if want.Bytes[i] != have.Bytes[i] {
				t.Fatalf("expected slice Bytes[%d] to be '%b' but got '%b'", i, want.Bytes[i], have.Bytes[i])
			}
		}
		if *have.Ptr != *want.Ptr {
			t.Fatalf("expected Ptr value to be '%v' but got '%v'", *have.Ptr, *want.Ptr)
		}
		if have.Struct != want.Struct {
			t.Fatalf("expected Ptr value to be '%v' but got '%v'", have.Struct, want.Struct)
		}
		if have.PtrIgnore != nil {
			t.Fatalf("expected PtrIgnore to be uninitialized but got %v", *have.PtrIgnore)
		}
		// t.Logf("have %v want %v", have, want)
		// t.Fail()
		
	})
	
}
func TestUnmarshalWithMapStringStringField(t *testing.T) {
	//SETUP
	t.Setenv("MAP_STRING_STRING", "k=v,k1='v1','k2'=v2,'k3'='v3',k4=v4")
	type Example struct {
		MAP_STRING_STRING map[string]string `env:"MAP_STRING_STRING"`
	}
	have := &Example{}
	want := &Example{
		MAP_STRING_STRING: map[string]string{
			"k": "v",
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
			"k4": "v4",
		},
	}

	if err := Unmarshal(  have); err != nil {
		t.Fatal(err)
	}
	
	for expected_key,expected_value := range want.MAP_STRING_STRING {
		if have.MAP_STRING_STRING[expected_key] != expected_value {
			t.Fatalf("expected MAP_STRING_STRING[%v] to be '%v' but was '%v'",expected_key, expected_value, have.MAP_STRING_STRING[expected_key],)
		}
	}
	
}

func TestUnmarshalWithMapStringIntField(t *testing.T) {
	//SETUP
	t.Setenv("MAP", "k=0,k1=1,'k2'='2','k3'='3',k4= 4  ")
	type Example struct {
		MAP map[string]int `env:"MAP"`
	}
	have := &Example{}
	want := &Example{
		MAP: map[string]int{
			"k": 0,
			"k1": 1,
			"k2": 2,
			"k3": 3,
			"k4": 4,
		},
	}

	if err := Unmarshal( have); err != nil {
		t.Fatal(err)
	}
	
	for expected_key,expected_value := range want.MAP {
		if have.MAP[expected_key] != expected_value {
			t.Fatalf("expected MAP[%v] to be '%v' but was '%v'",expected_key, expected_value, have.MAP[expected_key],)
		}
	}
	
}

func TestCustomReader(t *testing.T) {
	type Example struct {
		FOO string `env:"FOO;.env=./tag/env/foo.env,expand"`
		BAR string `env:"BAR;.env=./tag/env/foo.env,expand"`
		FOOBAR string `env:"FOOBAR;.env=./tag/env/foo.env,expand"`
		FOOBAZ string `env:"FOOBAZ;.env=./tag/env/foo.env,expand"`
	}
	have := &Example{}
	want := &Example{
		FOO: "FOO",
		BAR: "BAR",
		FOOBAR: "FOOBAR",
		FOOBAZ: "FOOBAZ", 
	}
	if err := Unmarshal(  have); err != nil {
		t.Fatal(err)
	}
	if *have != *want {
		t.Fatalf("expected %v to be %v", have, want)
	}
}

func TestUnmarshalStructWithUintSliceFields(t *testing.T) {
	//SETUP
	type BytesExample struct {
		SLICE_STRING []string `env:"SLICE_STRING;sep=','"`
		SLICE_BYTE []byte `env:"SLICE_BYTE;sep=''"`
		SLICE_UINT []uint `env:"SLICE_UINT;sep=' '"`
	}

	var test_cases = StructSliceTestCases[BytesExample] {
		{
			name: "String Slice Example",
			have: &BytesExample{},
			want: &BytesExample{
				SLICE_STRING: []string{"foo", "bar", "baz"},
				SLICE_BYTE: []byte("foobar"),
				SLICE_UINT: []uint{0,1,2},
			},
			setEnv: func(t *testing.T, tc NonComparableTestCase[BytesExample]) {
				t.Setenv("SLICE_STRING", "foo,bar,baz")
				t.Setenv("SLICE_BYTE", "foobar")
				t.Setenv("SLICE_UINT", "0 1 2")
			},
		},
		// {
		// 	name: "Uint Slice Syntax Error Example",
		// 	have: &BytesExample{},
		// 	want: &BytesExample{},
		// 	setEnv: func(t *testing.T, tc NonComparableTestCase[BytesExample]) {
		// 		t.Setenv("SLICE_UINT", "foobar")
		// 	},
		// 	expectedError: strconv.ErrSyntax,
		// },
		// {
		// 	name: "Uint Slice Error Example",
		// 	have: &BytesExample{},
		// 	want: &BytesExample{
		// 		SLICE_UINT: []uint{1234},
		// 	},
		// 	setEnv: func(t *testing.T, tc NonComparableTestCase[BytesExample]) {
		// 		t.Setenv("SLICE_UINT", "1234")
		// 	},
		// 	expectedError: strconv.ErrSyntax,
		// },
	}
	test_cases.Run(t, func(t *testing.T, have, want BytesExample) {
		if len(have.SLICE_STRING) != len(want.SLICE_STRING) {
			t.Fatalf("expected SLICE_STRING to have length '%d' but got '%d'", len(want.SLICE_STRING), len(have.SLICE_STRING))
		}
		for i := range want.SLICE_STRING {
			if want.SLICE_STRING[i] != have.SLICE_STRING[i] {
				t.Fatalf("expected SLICE_STRING[%d] to be '%s' but got '%s'", i, want.SLICE_STRING[i], have.SLICE_STRING[i])
			}
		}
		if len(have.SLICE_BYTE) != len(want.SLICE_BYTE) {
			t.Fatalf("expected SLICE_BYTE to have length '%d' but got '%d'", len(want.SLICE_BYTE), len(have.SLICE_BYTE))
		}
		for i := range want.SLICE_BYTE {
			if want.SLICE_BYTE[i] != have.SLICE_BYTE[i] {
				t.Fatalf("expected SLICE_BYTE[%d] to be '%b' but got '%b'", i, want.SLICE_BYTE[i], have.SLICE_BYTE[i])
			}
		}

		if len(have.SLICE_UINT) != len(want.SLICE_UINT) {
			t.Fatalf("expected SLICE_UINT to have length '%d' but got '%d'", len(want.SLICE_UINT), len(have.SLICE_UINT))
		}
		for i := range want.SLICE_UINT {
			if want.SLICE_UINT[i] != have.SLICE_UINT[i] {
				t.Fatalf("expected SLICE_UINT[%d] to be '%b' but got '%b'", i, want.SLICE_UINT[i], have.SLICE_UINT[i])
			}
		}
	})

	
}
type TextUnmarshalerExample struct {
		Name string
		Age int
}
func (t *TextUnmarshalerExample) UnmarshalText(data []byte) (err error) {
	parts := bytes.Split(data, []byte(","))
	t.Name = string(parts[0])
	
	return
}
func TestUnmarshalStructTextUnmarshalerFields(t *testing.T) {
	//SETUP

	
	type Example struct {
		SLICE_TEXT_UNMARSHALER []*TextUnmarshalerExample `env:"SLICE_STRING;sep=','"`
	}

	var test_cases = StructSliceTestCases[Example] {
		{
			name: "Empty TextUnmarshalerExample",
			have: &Example{},
			want: &Example{},
			setEnv: func(t *testing.T, tc NonComparableTestCase[Example]) {},
		},
		{
			name: "Empty TextUnmarshalerExample",
			have: &Example{},
			want: &Example{
				SLICE_TEXT_UNMARSHALER: []*TextUnmarshalerExample{
					{
						Name:"JOE",
					},
					{
						Name:"BOB",
					},
				},
			},
			setEnv: func(t *testing.T, tc NonComparableTestCase[Example]) {
				t.Setenv("SLICE_STRING", "JOE|1,BOB|2")
			},
		},
	}
	test_cases.Run(t, func(t *testing.T, have, want Example) {
		if len(have.SLICE_TEXT_UNMARSHALER) != len(want.SLICE_TEXT_UNMARSHALER) {
			t.Fatalf("expected SLICE_STRING to have length '%d' but got '%d' with value '%v'", len(want.SLICE_TEXT_UNMARSHALER), len(have.SLICE_TEXT_UNMARSHALER), have.SLICE_TEXT_UNMARSHALER)
		}
		// for i := range want.SLICE_STRING {
		// 	if want.SLICE_STRING[i] != have.SLICE_STRING[i] {
		// 		t.Fatalf("expected SLICE_STRING[%d] to be '%s' but got '%s'", i, want.SLICE_STRING[i], have.SLICE_STRING[i])
		// 	}
		// }
		
	})
	
}

// #################################### END TEST CASES ####################################


// #################################### TEST CASE RUNNERS ####################################

//Use this test case for struct types that satisfy the comparable type constraint
type ComparableTestCase[T any] struct {
		name string
		have *T
		want *T
		expectedAnyError bool
		expectedError error
		setEnv func (*testing.T, ComparableTestCase[T])
}

type StructTestCases[T comparable] []ComparableTestCase[T]

func (test_cases StructTestCases[T]) Run(_T *testing.T) {
	for _, tc := range test_cases {
		_T.Run(tc.name, func(t *testing.T) {
			if tc.setEnv != nil {
				tc.setEnv(t, tc)
			}
			if err := Unmarshal( tc.have); err != nil {
				if tc.expectedAnyError {
					slog.Debug("Expected any error")
					return
				} else if tc.expectedError != nil && errors.Is(err, tc.expectedError) {
					slog.Info("Matched expected error")
					return 
				}else if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
					t.Fatalf("expected err %s but got %s", tc.expectedError, err)
				} else {
					t.Fatal(err)
				}
			} else if tc.expectedAnyError {
				t.Fatalf("Expected any err but error was nil")
			} else if tc.expectedError != nil {
				t.Fatalf("expected err %s but was nil", tc.expectedError)
			}
			if *tc.have != *tc.want {
				t.Fatalf("[%s] expected %v but was %v",tc.name, *tc.want, *tc.have)
			}
		})
		
	}
}

//Use this test case for struct types that DO NOT satisfy the comparable type constraint
type NonComparableTestCase[T any] struct {
		name string
		have *T
		want *T
		expectedAnyError bool
		expectedError error
		setEnv func (*testing.T, NonComparableTestCase[T])
}

type StructSliceTestCases[T any] []NonComparableTestCase[T]

func (test_cases StructSliceTestCases[T]) Run(_T *testing.T, assert func(t *testing.T, have, want T)) {
	for _, tc := range test_cases {
		_T.Run(tc.name, func(t *testing.T) {
			if tc.setEnv != nil {
				tc.setEnv(t, tc)
			}
			if err := Unmarshal(tc.have); err != nil {
				if tc.expectedAnyError {
					slog.Debug("Expected any error")
					return
				} else if tc.expectedError != nil && errors.Is(err, tc.expectedError) {
					slog.Info("Matched expected error")
					return 
				}else if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
					t.Fatalf("expected err %s but got %s", tc.expectedError, err)
				} else {
					t.Fatal(err)
				}
			}else {
				assert(t, *tc.have, *tc.want)
			}			
		})
		
	}
}



func TestSimpleExample(t *testing.T) {
	t.Setenv("STR", "123")
	t.Setenv("INT", "123")
	t.Setenv("F32", "123.123")
	t.Setenv("F64", "123.123")
	t.Setenv("DURATION", "5s")
	type Example struct {
		Str string `env:"STR"`
		Int int `env:"INT"`
		Duration types.Duration `env:"DURATION"`
		// Float32 float32 `env:"F32"`
		// Float64 float64 `env:"F64"`
	}
	have := Example{}
	want := Example{
		Str: "123",
		Int: 123,
		Duration: types.Duration{time.Second * 5},
		// Float32: 123.123,
		// Float64: 123.123,
	}
	err := Unmarshal(&have)
	if err != nil {
		t.Fatal(err)
	}
	if have != want {
		t.Fatalf("expected %v but got %v", want, have)
	}
}

func TestArrayFields(t *testing.T) {
	t.Setenv("STR_ARRAY", "1,2,3")
	type Example struct {
		StrArray [3]string `env:"STR_ARRAY"`
	}
	have := Example{}
	want := Example{
		StrArray: [3]string{"1", "2", "3"},
	}
	err := Unmarshal( &have)
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range want.StrArray {
		if have.StrArray[i] != expected {
			t.Fatalf("Expected .StrSlice[%d] to be %q but was %q", i,expected,have.StrArray[i])
		}
	}

}

func TestArrayFieldErrorsWhenNotEnoughValues(t *testing.T) {
	t.Setenv("STR_ARRAY", "1,2,3,4")
	type Example struct {
		StrArray [3]string `env:"STR_ARRAY"`
	}
	have := Example{}
	err := Unmarshal(&have)
	expected_err := "cannot unpack 4 values into array of length 3" //\nreflect: array index out of range
	if err == nil {
		t.Fatalf("expected unmarshal to recover from panic and return an error")
	}else if err.Error() != expected_err {
		t.Fatalf("expected to recover from '%s', but was '%s'",expected_err, err)
	}
}

func TestStringSlice(t *testing.T) {
	t.Setenv("SLICE", "1,2,3")
	t.Setenv("ByteSliceWithEmptySeparator", "1,2,3")
	t.Setenv("DURATIONS", "1s,2s,3s")
	type Example struct {
		StrSlice []string `env:"SLICE"`
		IntSlice []int `env:"SLICE"`
		ByteSlice []byte `env:"SLICE"`
		ByteSliceWithEmptySeparator []byte `env:"ByteSliceWithEmptySeparator;sep='',other=true"`
		Uint8SliceWithSep []uint8 `env:"SLICE;sep=','"`
		Uint8Slice []uint8 `env:"SLICE"`
		// DurationSlice []Duration `env:"DURATIONS;sep=,"`
	}
	have := Example{}
	want := Example{
		StrSlice: []string{"1","2","3"},
		IntSlice: []int{1,2,3},
		ByteSlice: []byte{1,2,3},
		ByteSliceWithEmptySeparator: []byte("1,2,3"),
		Uint8Slice: []uint8{1,2,3},
		Uint8SliceWithSep: []uint8{1,2,3},

		// DurationSlice: []Duration{{time.Second},{time.Second*2}, {time.Second*3}},
	}
	err := Unmarshal(&have)
	if err != nil {
		t.Fatal(err)
	}
	// CHECK .StrSlice
	if len(want.StrSlice) != len(have.StrSlice) {
		t.Fatalf("expected .StrSlice length to be %d but was %d",len(want.StrSlice), len(have.StrSlice) )
	}
	for i, expected := range want.StrSlice {
		if have.StrSlice[i] != expected {
			t.Fatalf("Expected .StrSlice[%d] to be %q but was %q", i,expected,have.StrSlice[i])
		}
	}
	//VALIDATE .IntSlice
	if len(want.IntSlice) != len(have.IntSlice) {
		t.Fatalf("expected .IntSlice length to be %d but was %d",len(want.IntSlice), len(have.IntSlice) )
	}
	for i, expected := range want.IntSlice {
		if have.IntSlice[i] != expected {
			t.Fatalf("Expected .IntSlice[%d] to be %q but was %q", i,expected,have.IntSlice[i])
		}
	}
	//VALIDATE .ByteSlice
	if len(want.ByteSlice) != len(have.ByteSlice) {
		t.Fatalf("expected .BytesSlice length to be %d but was %d",len(want.ByteSlice), len(have.ByteSlice) )
	}
	for i, expected := range want.ByteSlice {
		if have.ByteSlice[i] != expected {
			t.Fatalf("Expected .ByteSlice[%d] to be %q but was %q", i,expected,have.ByteSlice[i])
		}
	}
	//VALIDATE .ByteSliceWithEmptySeparator
	if len(want.ByteSliceWithEmptySeparator) != len(have.ByteSliceWithEmptySeparator) {
		t.Fatalf("expected .ByteSliceWithEmptySeparator length to be %d but was %d",len(want.ByteSliceWithEmptySeparator), len(have.ByteSliceWithEmptySeparator) )
	}
	for i, expected := range want.ByteSliceWithEmptySeparator {
		if have.ByteSliceWithEmptySeparator[i] != expected {
			t.Fatalf("Expected .ByteSliceWithEmptySeparator[%d] to be %q but was %q", i,expected,have.ByteSliceWithEmptySeparator[i])
		}
	}
	//VALIDATE .Uint8Slice
	if len(want.Uint8Slice) != len(have.Uint8Slice) {
		t.Fatalf("expected .Uint8Slice length to be %d but was %d",len(want.Uint8Slice), len(have.Uint8Slice) )
	}
	for i, expected := range want.Uint8Slice {
		if have.Uint8Slice[i] != expected {
			t.Fatalf("Expected .Uint8Slice[%d] to be %q but was %q", i,expected,have.Uint8Slice[i])
		}
	}
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("expected %v but got %v", want, have)
	}
	// t.Fatalf("expected %v but got %v", want, have)
}


