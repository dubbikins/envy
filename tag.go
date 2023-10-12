package envy

import (
	"context"
	"encoding"
	"fmt"
	"os"
	"reflect"
	"regexp"
)

type Tag interface {
	TagMiddleware
	UnmarshalField(context.Context, reflect.StructField) error
}

func zeroValueUnmarshaller(ctx context.Context, field reflect.StructField) error {
	t, err := GetTagContext(ctx)
	if err != nil {
		return err
	}
	//Set the zero value for fields that can't be parsed from an empty string
	if t.Value == "" {
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			t.Value = "0"
		case reflect.Bool:
			t.Value = "false"
		}
	}
	return t.UnmarshalText(t.Bytes())
}

// unmarshalling order
// 1. env
// 2. default
// 3. options
// 4. matches
// 5. required
// 6. zero values (default unmarshaller)
func NewTag(value reflect.Value) Tag {
	return &tag{
		value: value,
		tagUnmarshallers: []Middleware{
			WithRequiredTag,
			WithMatchesTag,
			WithOptionsTag,
			WithDefaultTag,
			WithEnvTag,
		},
		customState:     map[string]interface{}{},
		tagUnmarshaller: TagHandlerFunc(zeroValueUnmarshaller),
	}
}

type tag struct {
	FieldType        string
	FieldName        string
	value            reflect.Value
	customState      map[string]interface{}
	index            int
	Name             string
	Default          string
	Value            string
	Raw              string
	Options          []string
	Required         bool
	Matcher          *regexp.Regexp
	IgnoreNil        bool
	textUnmarshaller TextUnmarshallable
	tagUnmarshallers []Middleware
	tagUnmarshaller  TagHandler
}

func (t *tag) UnmarshalText(text []byte) (err error) {
	return t.textUnmarshaller.UnmarshalText(text)
}
func (t *tag) GetState() map[string]interface{} {
	return t.customState
}
func (t *tag) GetStateValue(key string) interface{} {
	return t.customState[key]
}
func (t *tag) getChainedUnmarshallers() TagHandler {
	for len(t.tagUnmarshallers) > 0 {
		next := t.Pop()

		t.tagUnmarshaller = next(t.tagUnmarshaller)
	}
	return t.tagUnmarshaller
}
func (tag *tag) UnmarshalField(ctx context.Context, field reflect.StructField) (err error) {
	tag.FieldType = field.Type.Name()
	tag.FieldName = field.Name

	if !tag.value.IsValid() {
		return INVALID_FIELD_ERROR
	}
	ref := tag.value.Addr().Interface()
	if custom_text_unmarshaller, ok := ref.(encoding.TextUnmarshaler); ok {
		tag.useTextUnmarshaller(custom_text_unmarshaller)
	} else {
		switch tag.value.Kind() {
		case reflect.Ptr:
			tag.useTextUnmarshaller(_pointer(tag.value))
		case reflect.Struct:
			tag.useTextUnmarshaller(_struct(tag.value))
		case reflect.Slice:
			tag.useTextUnmarshaller(_slice(tag.value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tag.useTextUnmarshaller(_int(tag.value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tag.useTextUnmarshaller(_uint(tag.value))
		case reflect.String:
			tag.useTextUnmarshaller(_string(tag.value))
		case reflect.Bool:
			tag.useTextUnmarshaller(_boolean(tag.value))
		case reflect.Float32, reflect.Float64:
			tag.useTextUnmarshaller(_float(tag.value))
		}
	}
	if err = tag.getChainedUnmarshallers().UnmarshalField(WithTagContext(ctx, tag), field); err != nil {
		return
	}
	err = tag.UnmarshalText(tag.Bytes())
	return
}
func (t *tag) Bytes() []byte {
	if t.Value == "" {
		return []byte(t.Default)
	}
	return []byte(t.Value)
}

func (t *tag) Pop() Middleware {
	if len(t.tagUnmarshallers) == 0 {
		return nil
	}
	u := t.tagUnmarshallers[0]
	t.tagUnmarshallers = t.tagUnmarshallers[1:]
	return u
}
func (t *tag) Push(us ...Middleware) {
	t.tagUnmarshallers = append(us, t.tagUnmarshallers...)
}
func (t *tag) Contents() string {
	return fmt.Sprintf(`|*****Tag Details*****|
   StructFieldName: %s
   StructFieldType: %s
   Environment Variable Name: %s
   Aquired Value:%s
   Set Value: %s
   RawTag: %s
   Required: %t
`, t.FieldName, t.FieldType, t.Name, os.Getenv(t.Name), t.Value, t.Raw, t.Required)
}
func (t *tag) useTextUnmarshaller(u TextUnmarshallable) {
	t.textUnmarshaller = u
}
