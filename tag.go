package envy

import (
	"bytes"
	"context"
	"encoding"
	"fmt"
	"os"
	"reflect"
	"regexp"
)

type TagUnmarshaler interface {
	UnmarshalField(context.Context, reflect.StructField) error
}

func zeroValueUnmarshaller(ctx context.Context, field reflect.StructField) error {
	t, err := GetTagContext(ctx)
	if err != nil {
		return err
	}

	//Set the zero value for fields that can't be parsed from an empty string
	if t.Content == "" {
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			t.Content = "0"
		case reflect.Bool:
			t.Content = "false"
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
func NewTag(value reflect.Value, parent reflect.Value) (t *Tag, err error) {
	t = &Tag{
		Value:       value,
		Parent:      parent,
		customState: map[string]interface{}{},
		handler:     TagHandlerFunc(zeroValueUnmarshaller),
	}
	return
}

type Tag struct {
	FieldType   string
	FieldName   string
	Value       reflect.Value
	Parent      reflect.Value
	customState map[string]interface{}
	index       int
	Name        string
	Default     string
	Content     string
	Raw         string
	Options     []string
	Required    bool
	Matcher     *regexp.Regexp
	IgnoreNil   bool
	unmarshaler encoding.TextUnmarshaler
	middleware  []Middleware
	handler     TagHandler
	Skip        bool
	buffer      bytes.Buffer
}

func (t *Tag) Write(p []byte) (n int, err error) {
	// w.value = p
	// return len(p), nil
	return t.buffer.Write(p)
}

func (t *Tag) Read(p []byte) (n int, err error) {
	// w.value = p
	// return len(p), nil
	return t.buffer.Read(p)
}

func (t *Tag) UnmarshalText(text []byte) (err error) {
	return t.unmarshaler.UnmarshalText(text)
}
func (t *Tag) GetState() map[string]interface{} {
	return t.customState
}
func (t *Tag) GetStateValue(key string) interface{} {
	return t.customState[key]
}

func chainMiddleware(handler TagHandler, middlewares []Middleware) TagHandler {
	for len(middlewares) > 0 {
		next := middlewares[0]
		middlewares = middlewares[1:]
		handler = next(handler)
	}
	return handler
}
func (tag *Tag) UnmarshalField(ctx context.Context, field reflect.StructField) (err error) {
	tag.FieldType = field.Type.Name()
	tag.FieldName = field.Name
	if !tag.Value.IsValid() {
		return INVALID_FIELD_ERROR
	}
	ref := tag.Value.Addr().Interface()
	if custom_text_unmarshaller, ok := ref.(encoding.TextUnmarshaler); ok {
		tag.useTextUnmarshaller(custom_text_unmarshaller)
	} else {
		switch tag.Value.Kind() {
		case reflect.Ptr:
			tag.useTextUnmarshaller(_pointer(tag.Value))
		case reflect.Struct:
			tag.useTextUnmarshaller(_struct(tag.Value))
		case reflect.Slice:
			tag.useTextUnmarshaller(_slice(tag.Value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tag.useTextUnmarshaller(_int(tag.Value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tag.useTextUnmarshaller(_uint(tag.Value))
		case reflect.String:
			tag.useTextUnmarshaller(_string(tag.Value))
		case reflect.Bool:
			tag.useTextUnmarshaller(_boolean(tag.Value))
		case reflect.Float32, reflect.Float64:
			tag.useTextUnmarshaller(_float(tag.Value))
		default:
			//If the type is not one of these values, then it's likely an interface type and cannot be set
			//Simply return and ignore the values
			return
		}
	}
	var options *Options
	if options, err = GetOptionsContext(ctx); err != nil {
		return
	}
	if err = chainMiddleware(tag.handler, options.Middleware).UnmarshalField(WithTagContext(ctx, tag), field); err != nil {
		return
	}
	if !tag.Skip {
		err = tag.UnmarshalText(tag.Bytes())
	}
	return
}
func (t *Tag) Bytes() []byte {
	if t.Content == "" {
		return []byte(t.Default)
	}
	return []byte(t.Content)
}

func (t *Tag) DefaultMiddleware() []Middleware {
	return []Middleware{
		WithRequiredTag,
		WithMatchesTag,
		WithOptionsTag,
		WithEnvTag,
		WithDefaultTag,
	}
}

func (t *Tag) Push(us ...Middleware) {
	t.middleware = append(us, t.middleware...)
}
func (t *Tag) Contents() string {
	return fmt.Sprintf(`|*****Tag Details*****|
   StructFieldName: %s
   StructFieldType: %s
   Environment Variable Name: %s
   Aquired Value:%s
   Set Value: %s
   RawTag: %s
   Required: %t
`, t.FieldName, t.FieldType, t.Name, os.Getenv(t.Name), t.Content, t.Raw, t.Required)
}
func (t *Tag) useTextUnmarshaller(u encoding.TextUnmarshaler) {
	t.unmarshaler = u
}

type TagParser interface {
	TagName() string
	Handler() TagHandler
}
