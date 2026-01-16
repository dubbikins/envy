package tag

import (
	"errors"
	"fmt"

	"github.com/dubbikins/envy/v2/text"
)

// ################################ TAG LEXER START STATES ################################

func LexEnvironmentVariableTag(l text.Lexer[Token]) (next text.StateFn[Token]) {
	
	return LexEnvironmentVariableName
}

func LexTemplate(l text.Lexer[Token]) (next text.StateFn[Token]) {
	l.AcceptRunFn(IsAny)
	l.Emit(TokenIdent)
	
	return 
}


// ################################ LEXER STATES ################################ 

func EnvVarAllowedFirstChar (r rune) bool {
	return text.IsUpperOrLower(r) || r == '$' || r == '_'  
}

func EnvVarAllowedPostFirstChar (r rune) bool {
	return EnvVarAllowedFirstChar(r) || text.IsDigit(r) || r == '_' || r == '{' || r == '}' || r == '.'
}

func IsOptionKeyCharacter (r rune) bool {
	return r != ',' && r != '\'' && r != '='
}



func IsAny (r rune) bool {
	return true
}

func IsUpperAndLowerLettersDigitsAndUnderscore (r rune) bool {
	return text.IsUpperOrLower(r) ||  (r >= '0' && r <= '9') || r == '_'
}





func LexEnvironmentVariableName(l text.Lexer[Token]) (next text.StateFn[Token]) {
	l.AcceptRun(" ")
	l.Ignore()
	if l.AcceptRune('-'){
		l.Emit(TokenIdent)
		return nil
	}
	l.AcceptAnyOfFn(EnvVarAllowedFirstChar)
	l.AcceptRunFn(EnvVarAllowedPostFirstChar)
	if l.HasBufferedText() {
		l.Emit(TokenIdent)
	}
	l.AcceptRun(" ")
	l.Ignore()
	switch l.Peek() {
	case ';':
		return LexSemiColon
	case '|':
		return LexPipe
	case text.EOF :
		// l.Emit(EOF)
	default:
		
		return l.ErrorState(fmt.Errorf("unexpected token in input %s", l.RemainingText()))
	}
	return
}



func LexOption(l text.Lexer[Token]) (next text.StateFn[Token]) {
	if l.Peek() == text.EOF {
		return
	}
	//Lex the option key first
	return LexOptionKey
	
}

func LexOptionKey(l text.Lexer[Token]) (next text.StateFn[Token]) {
	l.AcceptRun(" ")
	l.Ignore()
	//Option key must start with an upper or lowercase letter
	if !l.AcceptAnyOfFn(IsOptionKeyCharacter) {
		return l.ErrorState(fmt.Errorf("LEX ERROR(LexOption): Option Key must start with upper and lowercase letter, have '%s' ",l.RemainingText()))
	}
	l.AcceptRunFn(IsOptionKeyCharacter)
	
	//Lex the separator
	accepted := l.AcceptRun(" ")
	var nextChar = l.Peek()
	var tokenType Token
	var noBufferedTextErr error
	// var atEOF bool
	switch nextChar {
	case '=':
		noBufferedTextErr = errors.New("LEX ERROR(LexOption): no option key provided before '=' in " + l.BufferedText())
		tokenType = TokenKey
		next =  LexOptionEquals
	case ',',text.EOF :
		noBufferedTextErr = errors.New("LEX ERROR(LexOption): no option flag provided before ',' in " + l.BufferedText())
		tokenType = TokenFlag
		next =  LexOption
	default:
		return l.ErrorState(errors.New("LEX ERROR(LexOption): unknown token [" + string(nextChar) + "] after option key in "+ l.BufferedText()))
	}	
	for accepted > 0 {
		l.Backup()
		accepted -= 1
	}
	l.AcceptRun(" ")
	l.Ignore()
	if l.HasBufferedText() {
		l.Emit(tokenType)
	}else {
		return l.ErrorState(noBufferedTextErr)
	}
	return 
}

// func LexOptionSeparator(l text.Lexer[Token]) (next text.StateFn[Token]) {
// 	l.AcceptRun(" ")
// 	l.Ignore()
	
// 	var nextChar = l.Peek()
// 	switch nextChar {
// 	case text.EOF :
// 		l.Emit(EOF)
// 		return 
// 	case '=':
// 		return LexOptionEquals
// 	case ',':
// 		if l.HasBufferedText() {
// 			l.Emit(TokenFlag)
// 			return LexOption
// 		}else {
// 			return l.ErrorState(errors.New("LEX ERROR(LexOption): no option flag provided before ',' in " + l.BufferedText()))
// 		}
// 	default:
// 		return l.ErrorState(errors.New("LEX ERROR(LexOption): unknown token [" + string(nextChar) + "] after option key in "+ l.BufferedText()))
// 	}	
// }

func LexComma(l text.Lexer[Token]) (next text.StateFn[Token]) {
	switch l.Next() {
	case ',':
		l.Emit(TokenComma)
		return LexOption
	case text.EOF:
		return
	default:
		return l.ErrorState(fmt.Errorf("LEX ERROR(LexComma): expected ',' but was %c", l.Peek()))
	}
}

func LexQuotedOptionValue(l text.Lexer[Token]) (next text.StateFn[Token]) {
	if !l.AcceptRune('\'') {
		return l.ErrorState(errors.New("LEX ERROR(LexQuotedOptionValue): expected single quote in " + l.BufferedText()))
	}
	l.AcceptRunFn(func(r rune) bool {
		return r != '\'' && r != text.EOF 
	})
	if !l.AcceptRune('\'') {
		return l.ErrorState(errors.New("LEX ERROR(LexQuotedOptionValue): expected single quote in " + l.BufferedText()))
	}
	l.Emit(TokenQuotedValue)
	return LexComma 
}

func LexOptionValue(l text.Lexer[Token]) (next text.StateFn[Token]) {
	l.AcceptRun(" ")
	l.Ignore()
	switch l.Peek() {
	case '\'' :
		return LexQuotedOptionValue
	case text.EOF:
		return l.ErrorState(errors.New("LEX ERROR(LexOptionValue): unexpected EOF after option key= in " + l.BufferedText()))
	}
	l.AcceptRunFn(func(r rune) bool {
		return r != ',' && r != text.EOF 
	})
	if !l.HasBufferedText() {
		return l.ErrorState(fmt.Errorf("LEX ERROR(LexOptionValue): unexpected ',' after %s. if you meant to specify a flag leave out the '='", l.BufferedText()))
	}
	l.Emit(TokenValue)
	return LexComma
}

func LexOptionEquals(l text.Lexer[Token]) (next text.StateFn[Token]) {
	if !l.AcceptRune('=') {
		return l.ErrorState(errors.New("LEX ERROR(LexOptionEquals): expected '=' after option key in " + l.BufferedText()))
	}
	l.Emit(TokenEquals)
	return LexOptionValue
}

func LexPipe(l text.Lexer[Token]) (next text.StateFn[Token]) {
	if !l.AcceptRune('|') {
		return l.ErrorState(fmt.Errorf("LEX ERROR(LexPipe): expected ';' but was %c", l.Peek()))
	}
	l.Emit(TokenPipe)
	return LexEnvironmentVariableName
}

func LexSemiColon(l text.Lexer[Token]) (next text.StateFn[Token]) {
	if !l.AcceptRune(';') {
		return l.ErrorState(fmt.Errorf("LEX ERROR(LexSemiColon): expected ';' but was %c", l.Peek()))
	}
	l.Emit(TokenSemiColon)
	return LexOption
}

