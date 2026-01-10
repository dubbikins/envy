package text

import (
	"fmt"
	"log/slog"

	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/text"
)

type Tag struct {
	Values []string
	Options map[string]string
	Flags map[string]struct{}
	Skipped bool
}

func Parse(tagName string, node *tag.Node, startState text.StateFn[Token]) ( err error) {
	if startState == nil {
		err = fmt.Errorf("parse lexing start state cannot be nil")
	}
	if node.Field() == nil {
		return 
	}
	slog.Debug("Parsing", "tag", tagName, "field", node.Field().Name)
	var tagValue, ok = node.Field().Tag.Lookup(tagName)
	if !ok   {
		return
	}

	var values = []string{}
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
	node.SetValueIterator( func(yield func(int, string) bool) {
		for i, value := range values {
			if !yield(i, value) {break}
		}
	})
	return
}

	


