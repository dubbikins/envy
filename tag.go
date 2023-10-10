package envy

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Tags []Taggable

type Taggable interface {
	ParseField(field reflect.StructField) error
}

const env_tagname = "env"
const required = "required"
const default_tag = "default"
const options_tagname = "options"
const matcher = "matcher"

type tag struct {
	Name      string
	Default   string
	Value     string
	Raw       string
	Options   []string
	Required  bool
	Matcher   string
	IgnoreNil bool
}

func getTag(field reflect.StructField) (t *tag, err error) {
	t = &tag{
		Raw: string(field.Tag.Get(env_tagname)),
	}
	if t.Raw == "" {
		return t, nil
	}
	parts := strings.SplitN(t.Raw, ";", 5)
	if len(parts) == 0 {
		return t, TAG_VALIDATION_ERROR
	} else if len(parts) > 5 {
		return t, INVALID_TAG_SYNTAX_ERROR
	}
	t.Name = parts[0]
	for _, part := range parts {
		if strings.HasPrefix(strings.TrimSpace(part), "default=") {
			t.Default = strings.Trim(strings.SplitN(part, "=", 2)[1], "\"")
		} else if strings.HasPrefix(strings.TrimSpace(part), "options=") {
			opts := strings.SplitN(part, "=", 2)[1]
			opts = strings.Trim(opts, "[]")
			t.Options = strings.Split(opts, ",")
		} else if strings.TrimSpace(part) == "required" {
			t.Required = true
		} else if strings.HasPrefix(strings.TrimSpace(part), "matcher=") {
			t.Matcher = strings.Trim(strings.SplitN(part, "=", 2)[1], "\"")
		} else if strings.ToLower(strings.TrimSpace(part)) == "ignorenil" {
			t.IgnoreNil = true
		}
	}

	t.Value = os.Getenv(t.Name)
	if t.Value == "" && t.Default != "" {
		t.Value = t.Default
	}

	details := fmt.Sprintf(`|******************************Tag Details******************************|
		StructField: %s
		OS Environment: %v
		Struct Tag Name: %s
		Environment Variable Name: %s
		Aquired Value:%s
		Set Value: %s
		RawTag: %s
		Required: %t
		`, field.Name, os.Environ(), env_tagname, t.Name, os.Getenv(t.Name), t.Value, t.Raw, t.Required)
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
	if t.Matcher != "" {
		matcher, err := regexp.Compile(t.Matcher)
		if err != nil {
			return t, err
		}

		if !matcher.Match([]byte(t.Value)) {
			return t, TagDoesMatchError(t.Name, t.Value, t.Matcher)
		}
	}

	if t.Value == "" && t.Default == "" && t.Required {
		return t, TagRequiredError(env_tagname, details)
	}
	//set the string version of the zero value if the value has no default or value set
	zero_value := ""
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		zero_value = "0"
	case reflect.Bool:
		zero_value = "false"
	}
	if t.Value == "" {
		t.Value = zero_value
	}

	return t, nil
}

type EnvTag struct {
}

func (tag *EnvTag) ParseField(field reflect.StructField) error {

	return nil
}

type OptionsTag struct {
	RawStructTag string
}

func (tag *OptionsTag) ParseField(field reflect.StructField) error {
	tag.RawStructTag = string(field.Tag.Get(options_tagname))
	return nil
}
