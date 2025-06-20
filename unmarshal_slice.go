package envy

import (
	"encoding"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
)

type _slice reflect.Value

func (s _slice) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		return
	}
	var num_elements int64
	slice_value := reflect.Value(s)
	temp_elem := reflect.MakeSlice(slice_value.Type(), 1, 1).Index(0)
	if temp_elem.Kind() == reflect.Interface  || temp_elem.Kind() == reflect.Ptr{
		slog.Warn("Cannot unmarshal slice of interface types or ptr", "slice_value", slice_value, "text", string(text))
		return
	}
	parts := strings.Split(string(text), ",")
	elem_kind := slice_value.Type().Elem().Kind()
	var is_custom bool

	slog.Info("Custom Unmarshaller check", "elem_kind", elem_kind, "slice_value", slice_value, "temp_elem", temp_elem, "text", string(text), "parts", parts)
	

	_, is_custom = temp_elem.Addr().Interface().(encoding.TextUnmarshaler)
	slog.Info("Custom Unmarshaller check", "is_custom", is_custom, "temp_elem", temp_elem, "text", string(text), "parts", parts)

	
	 // ensure the element type implements TextUnmarshaler
	if !is_custom && (elem_kind == reflect.Struct || elem_kind == reflect.Ptr || elem_kind == reflect.Slice)   {
		
		if len(text) == 0 {
			text = []byte("0")
		}
		//if the element type is a struct, pointer, or slice, we need to handle it differently
		//The value provided to default or returned from the env var provided in the env tag will be used to instantiate a slice of that size
		if num_elements, err = strconv.ParseInt(strings.ReplaceAll(string(text), ",", ""), 0, 0); err != nil {
			return err
		}
			slice_value.Set(reflect.MakeSlice(slice_value.Type(), int(num_elements), int(num_elements)))
		slog.Info("Unmarshalling slice of struct, ptr, or slices w/custom unmarshaller", "parts", parts, "len parts", len(parts), "slice_value", slice_value, "text", string(text),  "slice of type", slice_value.Type().Elem().Kind())

	}else {
			slog.Info("Unmarshalling slice", "parts", parts, "len parts", len(parts), "slice_value", slice_value, "text", string(text),  "slice of type", slice_value.Type().Elem().Kind())
			slice_value.Set(reflect.MakeSlice(slice_value.Type(), len(parts), len(parts)))
	}

	// 

	for i := 0; i < slice_value.Len(); i++ {
		var unmarshaler encoding.TextUnmarshaler
		var ok bool
		value := slice_value.Index(i)
	
		ref := value.Addr().Interface()
		if unmarshaler, ok = ref.(encoding.TextUnmarshaler); ok {
			slog.Info("Using custom unmarshaler", "value", value, "text", string(text), "parts", parts, "i", i, "ref", ref)
			if err = unmarshaler.UnmarshalText([]byte(strings.TrimSpace(parts[i]))); err != nil {
				return
			}
			
		} else {
		switch value.Kind() {
		//handle the recursive types
		case reflect.Ptr, reflect.Struct, reflect.Slice:
		// 	unmarshaler = _pointer(value)
		// case 
		// 	unmarshaler =  _struct(value)
		// case :
		// 	unmarshaler =  _slice(value)	
		
			if err = Unmarshal(ref); err != nil {
				return err
			}
		default:
			switch value.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if parts[i] == "" {
						parts[i] = "0" // default to false if empty
					}
					unmarshaler =  _int(value)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if parts[i] == "" {
						parts[i] = "0" // default to false if empty
					}
					unmarshaler = _uint(value)
				case reflect.String:
					unmarshaler = _string(value)
				case reflect.Bool:
					if parts[i] == "" {
						parts[i] = "false" // default to false if empty
					}
					unmarshaler = _boolean(value)
				case reflect.Float32, reflect.Float64:
					if parts[i] == "" {
						parts[i] = "0" // default to false if empty
					}
					unmarshaler = _float(value)
				default:
					//If the type is not one of these values, or the recursive types, then it's likely an interface type and cannot be set
					//Simply return and ignore the values
					return
				}
			if err= unmarshaler.UnmarshalText([]byte(strings.TrimSpace(parts[i]))); err != nil {
				slog.Error("Error unmarshalling slice element", "value", value, "text", string(text), "parts", parts, "i", i, "err", err)
				return err
			}
		}
	}
	
	}
	return nil
}


