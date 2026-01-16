package tag

import (
	"bufio"
	"bytes"
	"context"
	"encoding"
	"fmt"
	"iter"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/dubbikins/envy/v2/text"
)

type Node struct {
	parent *Node
	value reflect.Value
	field *reflect.StructField
	skip bool
	content bytes.Buffer
	options map[string]string
	flags map[string]struct{}
	//element_kind reflect.Kind
	separator byte
	position int
	iter iter.Seq2[int, string]
	// scanner bufio.Scanner
	ctx context.Context
	length int
}

// WithContext returns a shallow copy of n with its context changed
// to ctx. The provided ctx must be non-nil.
func (n *Node) WithContext(ctx context.Context) *Node {
	if ctx == nil {
		panic("nil context")
	}
	n2 := new(Node)
	*n2 = *n
	n2.ctx = ctx
	return n2
}

func (n *Node) SetTagValues(iter iter.Seq2[int, string]) {
	n.iter = iter
}

func (n *Node) Context() context.Context {
	return n.ctx
}

func (node *Node) Parse(tagName string,  startState text.StateFn[Token]) ( err error) {
	if startState == nil {
		err = fmt.Errorf("parse lexing start state cannot be nil")
	}
	var values = []string{}
	node.SetTagValues( func(yield func(int, string) bool) {
		for i, value := range values {
			if !yield(i, value) {break}
		}
	})
	if node.Field() == nil {
		return 
	}
	var tagValue, ok = node.Field().Tag.Lookup(tagName)
	if !ok   {
		return
	}

	
	// if node.lexer == nil {
	// 	node.lexer = 
	// }
	var tokenizer text.Tokenizer[Token] = text.NewLexer(tagValue, startState)
	var opts bool
	var token text.Token[Token]
	loop:
	for {
		token, err = tokenizer.NextToken()
		if err != nil  {
			return
		}else if token.ValueFrom(tagValue) == "-" {
			node.Skip()
			return
		}
		if token.Type == 0{
			break loop
		}
		if !opts {
			switch token.Type{
			case TokenIdent:
				values = append(values, token.ValueFrom(tagValue))
			case TokenSemiColon:
				opts = true
			case TokenPipe:
				continue
			case EOF:
				break
			default:
				err = fmt.Errorf("unexpected token type %s with value %s", token.Type, token.ValueFrom(tagValue))
				return
			}
		}else {
			switch token.Type{
			case TokenKey:
				var key text.Token[Token] = token
				var value text.Token[Token]
				if token, err = tokenizer.NextToken(); err != nil {
					return
				}else if  token.Type!= TokenEquals {
					err = fmt.Errorf("expected = after item key %s but have %s", key.ValueFrom(tagValue), token.ValueFrom(tagValue))
					return
				}
				if value, err = tokenizer.NextToken(); err != nil {
					return
				}else if  value.Type!= TokenValue && value.Type != TokenQuotedValue{
					err = fmt.Errorf("expected value after key= %s", value.ValueFrom(tagValue))
					return
				}
				if value.Type == TokenQuotedValue {
					value.Pos.Start +=1
					value.Pos.End -=1
				}
				node.SetOption(key.ValueFrom(tagValue), value.ValueFrom(tagValue))
				if token, err = tokenizer.NextToken(); err != nil {
					return
				}else if token.Type!= TokenComma && token.Type!= EOF{
					err = fmt.Errorf("expected , or EOF after key=value %s", tagValue[key.Start:token.End])
					return
				}
			case TokenFlag:
				node.SetFlag(token.ValueFrom(tagValue))
			case TokenComma:
				continue
			case EOF:
				break loop
			default:
				err = fmt.Errorf("exected key or EOF but was %s", token.ValueFrom(tagValue))
				return
			}
		}
		
	}

	return
}


func (n *Node) UnmarshalText(text []byte) (err error) {
	defer n.Reset()
	if n.value.CanAddr() {
		var unmarshaler encoding.TextUnmarshaler
		var implementsUnmarshaler bool
		if unmarshaler, implementsUnmarshaler = n.value.Addr().Interface().(encoding.TextUnmarshaler); implementsUnmarshaler {
			return unmarshaler.UnmarshalText(text)
		}
	}
	if len(text) == 0 {
		return
	}
	
	switch n.value.Kind() {
	case reflect.Pointer, reflect.Struct:
		
	case reflect.Map:
		if n.value.IsZero() {
			n.value.Set(reflect.MakeMap(n.value.Type()))
		}
		var sep, set = n.options["sep"]
		if len(sep) > 1 {
			return fmt.Errorf("invalid separator: limit 1 character")
		}
		if !set {
			n.separator = ','
			n.length = bytes.Count(text, []byte{n.separator})+1
		} else if len(sep) > 0 {
			n.separator = sep[0]
			n.length = bytes.Count(text, []byte{n.separator})+1
		}else {
			n.length = len(text)
		}
		
		for k, v := range n.map_kv_pairs(text) {
		
			var key = reflect.Indirect(reflect.New(reflect.TypeOf(n.value.Interface()).Key()))
			var val = reflect.Indirect(reflect.New(reflect.TypeOf(n.value.Interface()).Elem()))
			var _key = n.Descendant(key, nil)
			var _value = n.Descendant(val, nil)
			if err = _key.UnmarshalText(k); err != nil {
				return
			}
			if err = _value.UnmarshalText(v); err != nil {
				return
			}
			
			n.value.SetMapIndex(reflect.Indirect(_key.value), reflect.Indirect(_value.value))
		}
	case reflect.Slice,reflect.Array:
		text = bytes.Trim(text, "[{()}]")
		if len(text) == 0 {
			return
		}
		var sep, set = n.options["sep"]
		if len(sep) > 1 {
			return fmt.Errorf("invalid separator: limit 1 character")
		}
		if !set {
			n.separator = ','
			n.length = bytes.Count(text, []byte{n.separator})+1
		} else if len(sep) > 0 {
			n.separator = sep[0]
			n.length = bytes.Count(text, []byte{n.separator})+1
		}else {
			n.length = len(text)
		}
		
		if n.value.Kind() == reflect.Array && n.length > n.value.Len() {
			return fmt.Errorf("cannot unpack %d values into array of length %d", n.length, n.value.Len())
		}
		if n.value.Kind() == reflect.Slice && n.value.Len() < n.length {
			n.value.Grow(n.length)
			n.value.SetLen(n.length)
		}
		var i int
		for v := range n.slice_element_values(text) {
			var _value = n.Descendant(n.value.Index(i), nil)
			if err = _value.UnmarshalText(v); err != nil {
				return
			}
			i++
		}
	case reflect.String:
		//var b1 = 
		n.value.SetString(unsafe.String(&text[0], len(text)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64
		if val, err = strconv.ParseInt(string(text), 0, n.bitSize()); err == nil {
			n.value.SetInt(val)
		} 
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var val uint64
		//Since byte is an alias to the uint8 type, unmarshal on a type declared as []byte will end up getting unmarshalled
		//in this case. First we try to parse as uint8 from the text. If that fails, it was probably a slice of bytes, so 
		//we'll marshal the value as a rune
		if val, err = strconv.ParseUint(string(text), 0, n.bitSize()); err != nil {
			if n.bitSize() == 8 {
				n.value.SetUint(uint64(bytes.Runes(text)[0]))
				err = nil
			}
		} else {
			n.value.SetUint(val)
		}
	case reflect.Float32, reflect.Float64:
		var val float64
		if val, err = strconv.ParseFloat(string(text), n.bitSize()); err == nil {
			n.value.SetFloat(val)
		}
	
	case reflect.Complex64, reflect.Complex128:
		var val complex128
		if val, err = strconv.ParseComplex(string(text), n.bitSize()); err == nil {
			n.value.SetComplex(val)
		}
	case reflect.Bool:
		var val bool
		switch strings.ToLower(string(text)) {
		case "1","yes", "on", "TRUE", "T", "t", "True", "true": 
			val = true
		case "0", "no", "off", "FALSE", "F", "f", "False", "false":
			val = false
		}
		n.value.SetBool(val)
	}
	return 
}


func (n *Node) MarshalText() (text []byte, err error) {
	if n.value.CanAddr() {
		var marshaler encoding.TextMarshaler
		var implementsMarshaler bool
		if marshaler, implementsMarshaler = n.value.Addr().Interface().(encoding.TextMarshaler); implementsMarshaler {
			return marshaler.MarshalText()
		}
	}
	if len(text) == 0 {
		return
	}
	var format string = "%v"
	
	switch n.value.Kind() {
	case reflect.String:
		//var b1 = 
		n.value.SetString(unsafe.String(&text[0], len(text)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		format = "%d" 
	case reflect.Float32, reflect.Float64:
		format = "%f" 
	case reflect.Complex64, reflect.Complex128:
		//TODO
	case reflect.Bool:
		format = "%t"
	}
	return []byte(fmt.Sprintf(format, n.Value().Interface())), nil
}

func (n *Node) String() string {
	text, _ := n.MarshalText()
	return string(text)
}


func (n *Node) pointerDescendants(yield func(Node) bool) {
	if n.skip {
		return
	}else if n.value.IsNil() {
		n.value.Set(reflect.New(n.value.Type().Elem()))
	}
	var descendant = n.Descendant(n.value.Elem(), nil)
	if !yield(descendant){return}
}

func (n *Node) structDescendants(yield func(Node) bool) {
	var element = n.value.Type()
	for i := element.NumField()-1; i >= 0; i-- {
		var field = element.Field(i)
		var value = reflect.Value(n.value).Field(i)
		var descendant = n.Descendant(value,&field)
		if !field.IsExported() || !value.CanAddr(){
			continue
		}	
		if !yield(descendant){break}
	}
}

func (n *Node) sliceDescendants(yield func(Node) bool) {
	var i = 0
	for i < n.value.Len() {
		var descendant = n.Descendant(n.value.Index(i),nil)
		descendant.position = i
		if !yield(descendant){break}
		i++
	}
}

func (n *Node) mapDescendants(yield func(Node) bool) {
	var i = 1
	
	iter := n.value.MapRange()
	for iter.Next() {
		// iter.N
		var value_descendant = n.Descendant(iter.Value(), nil) //reflect.Indirect(reflect.New(reflect.TypeOf(n.value.Interface()).Elem())
		var key_descendant = n.Descendant(iter.Key(),nil) //n.descendant(reflect.Indirect(reflect.New(reflect.TypeOf(n.value.Interface()).Key())),n.field)
		key_descendant.position = i * -1
		value_descendant.position = i
		if !yield(value_descendant){break}
		if !yield(key_descendant){break}
		//n.value.SetMapIndex(key_descendant.value, value_descendant.value)
		i++
	}
}

func (n *Node) primitiveDescendants(yield func(Node) bool) {}

func (n *Node) descendants() iter.Seq[Node] {
	if n.value.CanAddr() {
		var implementsTextUnmarshaler bool
		if _, implementsTextUnmarshaler = n.value.Addr().Interface().(encoding.TextUnmarshaler); implementsTextUnmarshaler {
			return n.primitiveDescendants
		}
	}
	switch n.value.Kind() {
	case reflect.Struct:
		return n.structDescendants
	case reflect.Pointer:
		return n.pointerDescendants
	case reflect.Slice, reflect.Array:
		return n.sliceDescendants
	case reflect.Map:
		return n.mapDescendants
	}
	return n.primitiveDescendants
}

func (r *Node) Descendant(value reflect.Value, field *reflect.StructField) (Node) {
	return Node{parent: r, field: field, value: value}
}


func (n *Node) slice_element_values(text []byte) iter.Seq[[]byte] {
	if n.value.Kind() != reflect.Slice && n.value.Kind() != reflect.Array{
		panic("can't get element values on non-slice node")
	}
	var scanner = *bufio.NewScanner(bytes.NewReader(text))	
	scanner.Split(n.slice_element_scanner_splitfn)
	return func(yield func([]byte) bool) {
		for scanner.Scan() && yield(scanner.Bytes()) {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
		}
	}
}

func (n *Node) slice_element_scanner_splitfn(data []byte, atEOF bool) (advance int, token []byte, err error){
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if n.separator == 0 {
		advance = 1
		token = []byte(fmt.Sprintf("%v", data[0]))
		return
	}
	
	if i := bytes.IndexByte(data, n.separator); i >= 0 {
		advance = i + 1
		token = data[0:i]
		return 
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}


func (n *Node) map_kv_pairs(text []byte) iter.Seq2[[]byte,[]byte] {
	if n.value.Kind() != reflect.Map {
		panic("can't get kv pairs on non-map node")
	}
	var scanner = *bufio.NewScanner(bytes.NewReader(text))	
	scanner.Split(n.map_kv_scanner_splitfn)
	return func(yield func([]byte,[]byte) bool) {
		var key, value string
		for scanner.Scan() {
			key = scanner.Text()
			if !scanner.Scan() {
				if scanner.Err() != nil {
					panic(scanner.Err())
				}
				// panic("error scanning map kv pairs: expected value, but EOF")
			}
			value = scanner.Text()
			if !yield([]byte(key), []byte(value)){break}
		}
	}
}

func (n *Node) map_kv_scanner_splitfn(data []byte, atEOF bool) (advance int, token []byte, err error){
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if atEOF {
		advance = len(data)
		token = data
	}else if i := bytes.IndexAny(data, string([]byte{n.separator, byte('=')})); i >= 0 {
		advance = i + 1
		token = data[0:i]
	}
	if len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\''{
		token = token[1:len(token)-1]
	}
	token = bytes.TrimSpace(token)
	return 
}

func (n *Node) bitSize() (bitBase int) {
	switch n.value.Kind() {
		case reflect.Int:
		bitBase = 0
	case reflect.Int8, reflect.Uint8:
		bitBase = 8
	case reflect.Int16, reflect.Uint16:
		bitBase = 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		bitBase = 32
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64:
		bitBase = 64
	case reflect.Complex128:
		bitBase = 128
	}
	return
}

func (r *Node) Skip() {
	r.skip = true
}

func (r *Node) Skipped() bool {
	return r.skip
}

func (r *Node) Write(p []byte) (int,error) {
	return r.content.Write(p)
}

func (r *Node) Reset(){
	// r.SetValueIterator(nil)
	r.content.Reset()
}

func (r Node) Bytes() []byte {
	return r.content.Bytes()
}

func (r Node) Options() map[string]string {
	if r.options == nil {
		r.options = map[string]string{}
	}
	return r.options
}

func (r Node) Option(option string) (value string, set bool) {
	value, set = r.options[option]
	return 
}

func (r *Node) SetOption(op, v string) {
	if r.options == nil {
		r.options = map[string]string{}
	}
	r.options[op]= v
}

func (r Node) Value() reflect.Value {
	return r.value
}

func (n *Node) TagValues() iter.Seq2[int, string] {
	if n.iter == nil {
		return func(yield func(int, string) bool) {}
	}
	return n.iter
}

func (n *Node) SetFlag(flag string) {
	if n.flags == nil {
		n.flags = map[string]struct{}{}
	}
	n.flags[flag]=struct{}{}
}

func (n *Node) FlagIsSet(flag string) (set bool){
	_, set = n.flags[flag]
	return
}

func (n Node) Field() *reflect.StructField {
	return n.field
}

func (n *Node) Parent() (*Node) {
	return n.parent
}

func (n *Node) Root() (root *Node) {
	root = n
	for root.parent != nil {
		root = root.parent
	}
	return
}


func NewRootNode( s any) Node {
	return Node{
	
			parent: nil, 
			value: reflect.ValueOf(s), 
			field: nil, 
			content: *bytes.NewBuffer(nil),
			options: map[string]string{},
		}
}

