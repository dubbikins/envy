package etcd

import (
	"fmt"

	tag_text "github.com/dubbikins/envy/v2/tag/text"
	"github.com/dubbikins/envy/v2/text"
)


func LexEctdTag(l text.Lexer[tag_text.Token]) (next text.StateFn[tag_text.Token]) {
	return LexPathKey
}


// ################################ LEXER STATES ################################ 


func PathCharacters (r rune) bool {
	return text.IsUpperOrLower(r) || text.IsDigit(r) || r == '_' || r == '/' || r == '.' || r == '-'
}


func LexPathKey(l text.Lexer[tag_text.Token]) (next text.StateFn[tag_text.Token]) {
	l.AcceptRun(" ")
	l.Ignore()
	if l.AcceptRune('-'){
		l.Emit(tag_text.TokenIdent)
		return nil
	}
	l.AcceptRunFn(PathCharacters)
	if l.HasBufferedText() {
		l.Emit(tag_text.TokenIdent)
	}
	l.AcceptRun(" ")
	l.Ignore()
	switch l.Peek() {
	case ';':
		return tag_text.LexSemiColon
	case '|':
		return tag_text.LexPipe
	case text.EOF :
		// l.Emit(EOF)
	default:
		
		return l.ErrorState(fmt.Errorf("unexpected token in input %s", l.RemainingText()))
	}
	return
}


