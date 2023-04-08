package envy

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const tagname = "env"

type Tag struct {
	Name     string
	Default  string
	Value    string
	Raw      string
	Options  []string
	Required bool
}

func (field *FieldReflection) Tag() (t *Tag, err error) {
	t = &Tag{
		Raw: string(field.StructField.Tag.Get(tagname)),
	}
	parts := strings.Split(t.Raw, ";")
	if len(parts) == 0 {
		err = TAG_VALIDATION_ERROR
		return
	} else if len(parts) > 4 {
		err = INVALID_TAG_SYNTAX_ERROR
		return
	}
	for _, part := range parts {
		if strings.HasPrefix(strings.TrimSpace(part), "default=") {
			t.Default = strings.SplitN(part, "=", 2)[1]
		} else if strings.HasPrefix(strings.TrimSpace(part), "options=") {
			opts := strings.SplitN(part, "=", 2)[1]
			opts = strings.Trim(opts, "[]")
			t.Options = strings.Split(opts, ",")
		} else {
			t.Name = part
		}
	}
	t.Value = field.Reader.Getenv(t.Name)
	if t.Value == "" && t.Default != "" {
		t.Value = t.Default
	}
	if t.Options != nil {
		matches := false
		for _, option := range t.Options {
			if t.Value == option {
				matches = true
			}
		}
		if !matches {
			return t, TagInvalidOptionError(t.Name, t.Value, t.Options)
		}
	}
	if t.Value == "" && t.Default == "" && t.Required {
		return t, TagRequiredError(tagname, t.Value)
	}
	//set the string version of the zero value if the value has no default or value set
	zero_value := ""
	fmt.Println(field.Type.Kind())
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		zero_value = "0"
	case reflect.Bool:
		zero_value = "false"
	}
	if t.Value == "" {
		t.Value = zero_value
	}
	log.Printf("Tag Name: %s\n", t.Value)
	log.Printf("Tag Value: %s\n", t.Value)
	log.Printf("Tag Raw Value: %s\n", t.Raw)
	log.Printf("Tag Default: %s\n", t.Default)
	log.Printf("Tag Required: %t\n", t.Required)
	log.Printf("Tag Options: %v\n", t.Options)
	return t, nil
}

func UnmarshalString(tag *Tag, value ValueReflection) error {
	value.SetString(tag.Value)
	return nil
}

func UnmarshalInt(tag *Tag, value ValueReflection) (err error) {
	var bitBase int
	switch value.Kind {
	case reflect.Int:
		bitBase = 0
	case reflect.Int8:
		bitBase = 8
	case reflect.Int16:
		bitBase = 16
	case reflect.Int32:
		bitBase = 32
	default:
		bitBase = 64
	}
	if val, err := strconv.ParseInt(strings.ReplaceAll(tag.Value, ",", ""), 0, bitBase); err != nil {
		return err
	} else {
		value.SetInt(val)
	}
	return
}
func UnmarshalUint(tag *Tag, value ValueReflection) (err error) {
	var bitsize int
	switch value.Kind {
	case reflect.Uint:
		bitsize = 0
	case reflect.Uint8:
		bitsize = 8
	case reflect.Uint16:
		bitsize = 16
	case reflect.Uint32:
		bitsize = 32
	default:
		bitsize = 64
	}
	if val, err := strconv.ParseUint(strings.ReplaceAll(tag.Value, ",", ""), 0, bitsize); err != nil {
		return err
	} else {
		value.SetUint(val)
	}
	return
}

func UnmarshalBool(tag *Tag, value ValueReflection) (err error) {
	if val, err := strconv.ParseBool(tag.Value); err != nil {
		return err
	} else {
		value.SetBool(val)
	}
	return
}

func UnmarshalFloat(tag *Tag, value ValueReflection) (err error) {
	bitsize := 0
	switch value.Kind {
	case reflect.Float32:
		bitsize = 32
	case reflect.Float64:
		bitsize = 64
	}
	if val, err := strconv.ParseFloat(strings.ReplaceAll(tag.Value, ",", ""), bitsize); err == nil {
		value.SetFloat(val)
	}
	return
}
