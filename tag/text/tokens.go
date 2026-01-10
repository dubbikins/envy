package text

import "fmt"

type Token int32 

const (
	EOF Token = iota
	TokenPipe // pipe |
	TokenSemiColon  //semi colon ;
	TokenComma //comma ,
	TokenEquals //=
	TokenIdent //identifier
	TokenFlag
	TokenKey
	TokenValue
	TokenQuotedValue
	TokenText
	// TokenHyphen //-
)

func (i Token) String() string {
	switch i {
	// case TokenHyphen:
	// 	return "Hyphen"
	case EOF:
		return "EOF"
	case TokenPipe:
		return "PIPE"
	case TokenComma:
		return "COMMA"
	case TokenEquals:
		return "EQUALS"
	case TokenSemiColon:
		return "SEMICOLON"
	case TokenFlag:
		return "FLAG"
	case TokenKey:
		return "KEY"
	case TokenValue:
		return "VALUE"
	case TokenQuotedValue:
		return "QUOTED_VALUE"
	case TokenIdent:
		return "IDENTIFIER"
	default:
		return fmt.Sprintf("%d", i)
	}
}