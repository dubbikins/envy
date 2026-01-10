package text

import (
	"testing"

	"github.com/dubbikins/envy/v2/text"
)

// func TestLex(t *testing.T) {
// 	var input = "VAR|VAR2;k=v"
// 	var lexer = New(input, LexEnvironmentVariableTag)
// 	var expected_tokens = []string{"VAR", "|", "VAR2",";", "k", "=", "v", "EOF"}
// 	var item Item
// 	var err error
// 	var i = 0
// 	for i < len(expected_tokens) {
// 		item, err = lexer.NextItem()
// 		if err != nil  {
// 			t.Fail()
// 		}else if item.StringFrom(input) != expected_tokens[i] {
// 			t.Fatalf("expected token at position[%d] to be '%s' but was '%s'", i, expected_tokens[i], item.StringFrom(input))

// 		}
// 		i++
// 		if i == len(expected_tokens) && item.Type != ItemEOF {
// 			t.Fatalf("Expected ItemEOF")
// 		}
// 	}
// }

// func TestLexAccept(t *testing.T) {
// 	var input = "FOOBAR"
// 	var lexer = New(input, LexEnvironmentVariableTag)
// 	if !lexer.Accept("FOO") {
// 		t.Fail()
// 	}
// 	if lexer.Accept("BAR") {
// 		t.Fail()
// 	}
// 	if !lexer.Accept("FOO") {
// 		t.Fail()
// 	}
// 	if !lexer.Accept("FOO") {
// 		t.Fail()
// 	}
// 	if !lexer.Accept("BAR") {
// 		t.Fail()
// 	}
// 	if !lexer.Accept("BAR") {
// 		t.Fail()
// 	}
// 	if !lexer.Accept("BAR") {
// 		t.Fail()
// 	}
// }

func TestLex2(t *testing.T) {
	var input = "FOO|BAZ;key1='value',key2='value'"
	var lxr = text.NewLexer(input, LexEnvironmentVariableTag)
	var expected_tokens = []string{"FOO", "|", "BAZ",";", "key1", "=", "'value'",",", "key2", "=", "'value'", ""}
	var token text.Token[Token]
	var err error
	var i = 0
	for i < len(expected_tokens) {
		token, err = lxr.NextToken()
		if err != nil  {
			t.Fail()
		}else if token.ValueFrom(input) != expected_tokens[i] {
			t.Fatalf("expected token at position[%d] to be '%s' but was '%s'", i, expected_tokens[i], token.ValueFrom(input))

		}
		i++
		if i == len(expected_tokens) && token.Type != EOF{
			t.Fatalf("Expected EOF")
		}
	}
}
