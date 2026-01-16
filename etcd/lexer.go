package etcd

import (
	"fmt"

	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/text"
)


func LexEctdTag(l text.Lexer[tag.Token]) (next text.StateFn[tag.Token]) {
	return LexPathKey
}


// ################################ LEXER STATES ################################ 


func PathCharacters (r rune) bool {
	return text.IsUpperOrLower(r) || text.IsDigit(r) || r == '_' || r == '/' || r == '.' || r == '-'
}


func LexPathKey(l text.Lexer[tag.Token]) (next text.StateFn[tag.Token]) {
	l.AcceptRun(" ")
	l.Ignore()
	if l.AcceptRune('-'){
		l.Emit(tag.TokenIdent)
		return nil
	}
	l.AcceptRunFn(PathCharacters)
	if l.HasBufferedText() {
		l.Emit(tag.TokenIdent)
	}
	l.AcceptRun(" ")
	l.Ignore()
	switch l.Peek() {
	case ';':
		return tag.LexSemiColon
	case '|':
		return tag.LexPipe
	case text.EOF :
		// l.Emit(EOF)
	default:
		
		return l.ErrorState(fmt.Errorf("unexpected token in input %s", l.RemainingText()))
	}
	return
}


