package envy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
)

type testStruct struct {
	unexported   string  `env:"TEST_ENV_UNEXPORTED"`
	String       string  `env:"TEST_ENV_STR"`
	StringPtr    string  `env:"TEST_ENV_STR_PTR"`
	IntPtr       int     `env:"TEST_ENV_INT_PTR"`
	IntMin       int     `env:"TEST_ENV_INT_MIN"`
	Int8Min      int8    `env:"TEST_ENV_INT8_MIN"`
	Int16Min     int16   `env:"TEST_ENV_INT16_MIN"`
	Int32Min     int32   `env:"TEST_ENV_INT32_MIN"`
	Int64Min     int64   `env:"TEST_ENV_INT64_MIN"`
	IntMax       int     `env:"TEST_ENV_INT_MAX"`
	Int8Max      int8    `env:"TEST_ENV_INT8_MAX"`
	Int16Max     int16   `env:"TEST_ENV_INT16_MAX"`
	Int32Max     int32   `env:"TEST_ENV_INT32_MAX"`
	Int64Max     int64   `env:"TEST_ENV_INT64_MAX"`
	UintMin      uint    `env:"TEST_ENV_UINT_MIN"`
	Uint8Min     uint8   `env:"TEST_ENV_UINT8_MIN"`
	Uint16Min    uint16  `env:"TEST_ENV_UINT16_MIN"`
	Uint32Min    uint32  `env:"TEST_ENV_UINT32_MIN"`
	Uint64Min    uint64  `env:"TEST_ENV_UINT64_MIN"`
	UintMax      uint    `env:"TEST_ENV_UINT_MAX"`
	Uint8Max     uint8   `env:"TEST_ENV_UINT8_MAX"`
	Uint16Max    uint16  `env:"TEST_ENV_UINT16_MAX"`
	Uint32Max    uint32  `env:"TEST_ENV_UINT32_MAX"`
	Uint64Max    uint64  `env:"TEST_ENV_UINT64_MAX"`
	Float32Min   float32 `env:"TEST_ENV_FLOAT32_MIN"`
	Float64Min   float64 `env:"TEST_ENV_FLOAT64_MIN"`
	Float32Max   float32 `env:"TEST_ENV_FLOAT32_MAX"`
	Float64Max   float64 `env:"TEST_ENV_FLOAT64_MAX"`
	BoolTrue     bool    `env:"TEST_ENV_BOOL_TRUE"`
	BoolFalse    bool    `env:"TEST_ENV_BOOL_False"`
	BoolPtrTrue  bool    `env:"TEST_ENV_BOOL_PTR_TRUE"`
	BoolPtrFalse bool    `env:"TEST_ENV_BOOL_PTR_FALSE"`
	BoolYes      bool    `env:"TEST_ENV_BOOL_YES"`
	BoolNo       bool    `env:"TEST_ENV_BOOL_NO"`
	BoolOn       bool    `env:"TEST_ENV_BOOL_ON"`
	BoolOff      bool    `env:"TEST_ENV_BOOL_OFF"`
	Bool1        bool    `env:"TEST_ENV_BOOL_1"`
	Bool0        bool    `env:"TEST_ENV_BOOL_0"`
	NestedStruct struct {
		Field string `env:"TEST_STRUCT_FIELD"`
	}
	NestedStructPointer *struct {
		Field string `env:"TEST_STRUCT_FIELD_PTR"`
	}
}

type TestStructWithPointer struct {
	NestedStructPointer *struct {
		Field string `env:"TEST_STRUCT_FIELD_PTR"`
	}
}
type testOutOfBoundsStruct struct {
	IntBelowMin     int     `env:"TEST_ENV_INT_BELOW_MIN"`
	Int8BelowMin    int8    `env:"TEST_ENV_INT8_BELOW_MIN"`
	Int16BelowMin   int16   `env:"TEST_ENV_INT16_BELOW_MIN"`
	Int32BelowMin   int32   `env:"TEST_ENV_INT32_BELOW_MIN"`
	Int64BelowMin   int64   `env:"TEST_ENV_INT64_BELOW_MIN"`
	IntMaxAboveMax  int     `env:"TEST_ENV_INT_ABOVE_MAX"`
	Int8AboveMax    int8    `env:"TEST_ENV_INT8_ABOVE_MAX"`
	Int16AboveMax   int16   `env:"TEST_ENV_INT16_ABOVE_MAX"`
	Int32AboveMax   int32   `env:"TEST_ENV_INT32_ABOVE_MAX"`
	Int64AboveMax   int64   `env:"TEST_ENV_INT64_ABOVE_MAX"`
	UintBelowMin    uint    `env:"TEST_ENV_UINT_BELOW_MIN"`
	Uint8BelowMin   uint8   `env:"TEST_ENV_UINT8_BELOW_MIN"`
	Uint16BelowMin  uint16  `env:"TEST_ENV_UINT16_BELOW_MIN"`
	Uint32BelowMin  uint32  `env:"TEST_ENV_UINT32_BELOW_MIN"`
	Uint64BelowMin  uint64  `env:"TEST_ENV_UINT64_BELOW_MIN"`
	UintAboveMax    uint    `env:"TEST_ENV_UINT_ABOVE_MAX"`
	Uint8AboveMax   uint8   `env:"TEST_ENV_UINT8_ABOVE_MAX"`
	Uint16AboveMax  uint16  `env:"TEST_ENV_UINT16_ABOVE_MAX"`
	Uint32AboveMax  uint32  `env:"TEST_ENV_UINT32_ABOVE_MAX"`
	Uint64AboveMax  uint64  `env:"TEST_ENV_UINT64_ABOVE_MAX"`
	Float32BelowMin float32 `env:"TEST_ENV_FLOAT32_BELOW_MIN"`
	Float64BelowMin float64 `env:"TEST_ENV_FLOAT64_BELOW_MIN"`
	Float32AboveMax float32 `env:"TEST_ENV_FLOAT32_ABOVE_MAX"`
	Float64AboveMax float64 `env:"TEST_ENV_FLOAT64_ABOVE_MAX"`
}

type defaultTestStruct struct {
	String       string `env:"TEST_ENV_STR;default=passed"`
	Int          int    `env:"TEST_ENV_INT;default=200"`
	Int8         int8   `env:"TEST_ENV_INT8;default=200"`
	Int16        int16  `env:"TEST_ENV_INT16;default=200"`
	Int32        int32  `env:"TEST_ENV_INT32;default=200"`
	Int64        int64  `env:"TEST_ENV_INT64;default=200"`
	Uint         int    `env:"TEST_ENV_UINT;default=200"`
	Uint8        uint8  `env:"TEST_ENV_UINT8;default=200"`
	Uint16       uint16 `env:"TEST_ENV_UINT16;default=200"`
	Uint32       uint32 `env:"TEST_ENV_UINT32;default=200"`
	Uint64       uint64 `env:"TEST_ENV_UINT64;default=200"`
	Bool         bool   `env:"TEST_ENV_BOOL;default=true"`
	NestedStruct struct {
		Field string `env:"TEST_STRUCT_FIELD;default=passed"`
	}
}

type optionsTestStruct struct {
	String string `env:"TEST_ENV_STR" options:"a,b,c"`
	Int    int    `env:"TEST_ENV_INT" options:"[-1,2,3]"`
	Int8   int8   `env:"TEST_ENV_INT8" options:"{-1,2,3}"`
	Int16  int16  `env:"TEST_ENV_INT16" options:"(-1,2,3)"`
	Int32  int32  `env:"TEST_ENV_INT32" options:"[-1,2,3]"`
	Int64  int64  `env:"TEST_ENV_INT64" options:"[-1,2,3]"`
	Uint   uint   `env:"TEST_ENV_UINT" options:"[1,2,3]"`
	Uint8  uint8  `env:"TEST_ENV_UINT8" options:"[1,2,3]"`
	Uint16 uint16 `env:"TEST_ENV_UINT16" options:"[1,2,3]"`
	Uint32 uint32 `env:"TEST_ENV_UINT32" options:"[1,2,3]"`
	Uint64 uint64 `env:"TEST_ENV_UINT64" options:"[1,2,3]"`
	Bool   bool   `env:"TEST_ENV_BOOL" options:"[yes,no]"`
}

type requiredTestStruct struct {
	String  string  `env:"TEST_ENV_STR" required:"true"`
	Int     int     `env:"TEST_ENV_INT" required:"true"`
	Int8    int     `env:"TEST_ENV_INT8" required:"true"`
	Int16   int     `env:"TEST_ENV_INT16" required:"true"`
	Int32   int     `env:"TEST_ENV_INT32" required:"true"`
	Int64   int     `env:"TEST_ENV_INT64" required:"true"`
	Uint    int     `env:"TEST_ENV_UINT" required:"true"`
	Uint8   int     `env:"TEST_ENV_UINT8" required:"true"`
	Uint16  int     `env:"TEST_ENV_UINT16" required:"true"`
	Uint32  int     `env:"TEST_ENV_UINT32" required:"true"`
	Uint64  int     `env:"TEST_ENV_UINT64" required:"true"`
	Float32 float32 `env:"TEST_ENV_FLOAT32" required:"true"`
	Float64 float64 `env:"TEST_ENV_FLOAT64" required:"true"`
	Bool    bool    `env:"TEST_ENV_BOOL" required:"true"`
}

type matches struct {
	String  string  `env:"TEST_ENV_STR" matches:"(.*).txt"`
	Int     int     `env:"TEST_ENV_INT" required:"true"`
	Int8    int     `env:"TEST_ENV_INT8" required:"true"`
	Int16   int     `env:"TEST_ENV_INT16" required:"true"`
	Int32   int     `env:"TEST_ENV_INT32" required:"true"`
	Int64   int     `env:"TEST_ENV_INT64" required:"true"`
	Uint    int     `env:"TEST_ENV_UINT" required:"true"`
	Uint8   int     `env:"TEST_ENV_UINT8" required:"true"`
	Uint16  int     `env:"TEST_ENV_UINT16" required:"true"`
	Uint32  int     `env:"TEST_ENV_UINT32" required:"true"`
	Uint64  int     `env:"TEST_ENV_UINT64" required:"true"`
	Float32 float32 `env:"TEST_ENV_FLOAT32" required:"true"`
	Float64 float64 `env:"TEST_ENV_FLOAT64" required:"true"`
	Bool    bool    `env:"TEST_ENV_BOOL" required:"true"`
}

var testingTKey = contextKey("__testing_T__")
var structDefKey = contextKey("__struct_defintion__")
var testTypeKey = contextKey("__test_type__")

func structWithFields(ctx context.Context, _type string) (context.Context, error) {
	switch _type {
	case "base":
		ctx = context.WithValue(ctx, structDefKey, &testStruct{})
	case "options":
		ctx = context.WithValue(ctx, structDefKey, &optionsTestStruct{})
	case "default":
		ctx = context.WithValue(ctx, structDefKey, &defaultTestStruct{})
	case "required":
		ctx = context.WithValue(ctx, structDefKey, &requiredTestStruct{})
	default:
		return ctx, errors.New("invalid type")
	}
	return context.WithValue(ctx, testTypeKey, _type), nil
}

func theEnvVarIsSetTo(ctx context.Context, key, value string) (context.Context, error) {

	t, ok := ctx.Value(testingTKey).(*testing.T)
	if !ok {
		return ctx, errors.New("testing.T not found in context")
	}
	t.Setenv(key, value)
	return ctx, nil
}

func callUnmarshal(ctx context.Context) (context.Context, error) {
	return ctx, Unmarshal(ctx.Value(structDefKey).(*testStruct))
}

func hasValues(ctx context.Context, expectedValues *godog.DocString) (context.Context, error) {
	_type, ok := ctx.Value(testTypeKey).(string)
	if !ok {
		return ctx, errors.New("struct test type not found in context")
	}
	var err error
	var actual any
	var expected any
	switch _type {
	case "base":
		actual = ctx.Value(structDefKey).(*testStruct)
		expected = &testStruct{}
	case "options":
		actual = ctx.Value(structDefKey).(*optionsTestStruct)
		expected = &optionsTestStruct{}
	case "default":
		actual = ctx.Value(structDefKey).(*defaultTestStruct)
		expected = &defaultTestStruct{}
	case "required":
		actual = ctx.Value(structDefKey).(*requiredTestStruct)
		expected = &requiredTestStruct{}
	default:
		return ctx, errors.New("invalid type")
	}
	err = json.Unmarshal([]byte(expectedValues.Content), expected)
	if err != nil {
		return ctx, err
	}
	if !cmp.Equal(&actual, &expected, cmp.AllowUnexported(testStruct{})) {
		comp := &Comparision{&actual, &expected}
		comp.ShowDeepUnequal()
		return ctx, fmt.Errorf("struct values do not match: actual(%v) != expected(%v)", actual, expected)
	}
	return ctx, nil
}

type Comparision struct {
	Actual   interface{}
	Expected interface{}
}

func Comp(actual, expected interface{}) *Comparision {
	return &Comparision{actual, expected}
}

func (c *Comparision) ShowDeepUnequal() {
	actual_value := reflect.ValueOf(c.Actual)
	expected_value := reflect.ValueOf(c.Expected)
	if actual_value.Kind() != expected_value.Kind() {
		fmt.Println("kind mismatch")
		return
	}
	_type := actual_value.Type()
	switch _type.Kind() {
	case reflect.Struct:
		fmt.Printf("Comparing Structs[%s]\n", _type.Name())
		for i := 0; i < _type.NumField(); i++ {
			if fmt.Sprintf("%s", actual_value.Field(i)) != fmt.Sprintf("%s", expected_value.Field(i)) {
				fmt.Printf("FieldName: %s\n\tActual: %s\n\tExpected %s\n", _type.Field(i).Name, actual_value.Field(i), expected_value.Field(i))
			}
		}
	case reflect.Pointer:
		fmt.Println("Comparing Pointers; dereferencing...")
		Comp(actual_value.Elem().Interface(), expected_value.Elem().Interface()).ShowDeepUnequal()
	default:
		fmt.Println("Comparing Primitives or collections")
		fmt.Printf("actual(%s) vs expected (%s)\n", actual_value, expected_value)

		fmt.Println(reflect.DeepEqual(c.Actual, c.Expected))
	}

}

func hasNoErrors(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a "(base|default|options|required)" struct test struct is instantiated$`, structWithFields)
	ctx.Step(`^the environment variable "(.+)" is set to "(.+)"$`, theEnvVarIsSetTo)
	ctx.Step(`^the struct is passed as an argument to Unmarshal$`, callUnmarshal)
	ctx.Step(`^the struct should have the following values:$`, hasValues)
	ctx.Step(`^there should be no errors$`, hasNoErrors)
}
func TestingTScrenarioProvider(t *testing.T) func(*godog.ScenarioContext) {
	return func(gctx *godog.ScenarioContext) {
		gctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			os.Clearenv()
			return context.WithValue(ctx, testingTKey, t), nil
		})
		gctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
			os.Clearenv()
			return ctx, nil
		})
		InitializeScenario(gctx)
	}

}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: TestingTScrenarioProvider(t),
		Options: &godog.Options{
			StopOnFailure: true,
			Strict:        true,
			Format:        "progress",
			Paths:         []string{"tests/features"},
			TestingT:      t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
