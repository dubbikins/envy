package envy

import (
	"encoding"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

type Options struct {
	Middleware []Middleware
}

type opt struct {
	field reflect.StructField
	rvalue reflect.Value
	literal string
}

type tagUnmarshallerOptions struct {
	OverrideValues bool 
	_opts_cache map[string]opt
}

func (opts *tagUnmarshallerOptions) Load(value string) error {
	if !strings.HasPrefix(string(value), "@") {
		return errors.New("when loading unmarshallerOptions the values must start with @")
	}else {
		slog.Debug("loading tagUnmarshallerOptions", "value", value)
	}
	var option_name string 
	var option_value string
	var segments = strings.SplitN(string(value), "=", 2)
	if len(segments) > 0 {
		option_name = segments[0]
	}
	if len(segments) > 1 {
		option_value = segments[1]
	}
	if strings.Contains(string(value), "=") {
		//cache the options fields so successive calls to Load don't have to use reflection
		//to find the field
		//this is a performance optimization
		if opts._opts_cache == nil {
			opts._opts_cache = make(map[string]opt)
			rvalue := reflect.ValueOf(opts)
			elem := rvalue.Elem()
		
			for i := 0; i < elem.NumField(); i++ {
				field := elem.Type().Field(i)
				field_value := rvalue.Elem().Field(i)
				if !field.IsExported() || !field_value.CanSet() {
					continue
				}
				opts._opts_cache[option_name] = opt{
					field: field,
					rvalue: field_value,
					literal: option_value,
				}
			}
		}

		// unmarshal the field according to the type
		if opt, ok := opts._opts_cache[option_name] ; ok {
			
			var unmarshaller encoding.TextUnmarshaler
			switch opt.field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				unmarshaller = _int(opt.rvalue)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				unmarshaller = _uint(opt.rvalue)
			case reflect.String:
				unmarshaller = _string(opt.rvalue)
			case reflect.Bool:
				unmarshaller = _boolean(opt.rvalue)
			case reflect.Float32, reflect.Float64:
				unmarshaller = _float(opt.rvalue)
			default:
				panic(fmt.Sprintf("unmarshallerOptions: unknown option field type %s", opt.field.Type))
			}
			if unmarshaller != nil {
				if err := unmarshaller.UnmarshalText([]byte(opt.literal)); err != nil {
					return fmt.Errorf("error unmarshalling option %s: %w", option_name, err)
				}
			}
		}
	}
	return nil
}




func (opts *Options) Unmarshal(s any) {

}

func DefaultMiddleware() []Middleware {
	return []Middleware{
		WithRequiredTag,
		WithMatchesTag,
		WithOptionsTag,
		WithEnvTag,
		WithDefaultTag,
		WithEnvyGlobalTag,
	}
}

var default_options *Options

func init() {
	default_options = &Options{
		Middleware: DefaultMiddleware(),
	}
}
