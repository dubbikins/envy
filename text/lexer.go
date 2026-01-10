package text

import (
	"log/slog"
	"strings"
	"unicode/utf8"
)

type Tokenizer[T TokenType] interface {
	NextToken() (next Token[T], err error)
}


type TokenType interface {
	~int8  | ~int16 | ~int | ~int32 | ~int64 
	String() string
}

type StateFn[T TokenType] func(Lexer[T]) StateFn[T]


type Pos struct {
	Start int
	End int
}

type Token[T TokenType] struct {
	Type T
	Pos
}

func (t Token[T]) ValueFrom(src string) string {
	return src[t.Start:t.End]
}

type Lexer[T TokenType] interface {
	Emit(tokenType T)
	Next() (char rune)
	Peek() (char rune)
	Backup()
	Ignore()
	Reset(startState StateFn[T])
	ErrorState(cause error) StateFn[T]
	Accept(valid string) (accepted bool)
	AcceptRun(valid string) (accepted int) 
	AcceptAnyOfFn(validator func(rune) bool) (accepted bool)
	AcceptRunFn(validator func(rune) bool) (accepted bool)
	AcceptRune(rune rune) (accepted bool)
	BufferedText() string 
	RemainingText() string
	HasBufferedText() bool
	emitError(err error)
}

type lexer[T TokenType] struct {
	input string
	// start int
	// pos int
	pos Pos
	width int
	state StateFn[T]
	items chan Token[T]
	errors chan error
}

func NewLexer[T TokenType](input string, initialState StateFn[T]) *lexer[T] {
	return &lexer[T]{
		input: input,
		state: initialState,
		items: make(chan Token[T], 2),
		errors: make(chan error, 1),
	}
}

func (lxr *lexer[T]) NextToken() (next Token[T], err error) {
	for {
		select {
		case next = <-lxr.items:
			return
		case err = <- lxr.errors:
			return
		default:
			if lxr.state == nil {
				return
			}
			lxr.state = lxr.state(lxr)
		}
	}
}

func (lxr *lexer[T]) Emit(tokenType T) {
	slog.Debug("emitting", "token", lxr.input[lxr.pos.Start:lxr.pos.End], "type", tokenType, "pos", lxr.pos)
	lxr.items <- Token[T]{Type: tokenType, Pos: lxr.pos}
	lxr.pos.Start = lxr.pos.End
}

const EOF rune = 0

func (lxr *lexer[T]) Next() (char rune) {
	if lxr.pos.End >= len(lxr.input) {
		lxr.width = 0
		return EOF
	}
	char, lxr.width = utf8.DecodeRuneInString(lxr.input[lxr.pos.End:])
	lxr.pos.End += lxr.width
	return
}

func (lxr *lexer[T]) Peek() (char rune) {
	char = lxr.Next()
	lxr.Backup()
	return
}
func (lxr *lexer[T]) Backup() {
	lxr.pos.End -= lxr.width
}

func (lxr lexer[T]) Ignore() {
	lxr.pos.Start = lxr.pos.End
}

func (lxr *lexer[T]) emitError(err error) {
	lxr.errors <- err
}


func (lexer[T]) ErrorState(cause error) StateFn[T] {
	if cause == nil {
		panic("StopOnErrorState called with nil error")
	}
	return func(lxr Lexer[T]) StateFn[T] {
		lxr.emitError(cause)
		return nil
	}
}


func (lxr lexer[T]) Accept(valid string) (accepted bool) {
	if strings.ContainsRune(valid, lxr.Next()) {
		return true
	}
	lxr.Backup()
	return
}


func (lxr *lexer[T]) AcceptRun(valid string) (accepted int) { //
	for strings.ContainsRune(valid, lxr.Next()) {
		accepted += 1
	}
	lxr.Backup()
	return
}

func (lxr *lexer[T]) Reset(startState StateFn[T]) {
	if lxr.items != nil {
		close(lxr.items)
	}
	if lxr.errors != nil {
		close(lxr.errors)
	}
	lxr.pos.Start = 0
	lxr.pos.End = 0
	lxr.items = make(chan Token[T], 2)
	lxr.errors = make(chan error, 1)
	lxr.state = startState
}

func (lxr *lexer[T]) AcceptAnyOfFn(validator func(rune) bool) (accepted bool) {
	var next = lxr.Next()
	if accepted = validator(next); accepted {
		return
	}
	lxr.Backup()
	return
}

func (lxr *lexer[T]) AcceptRunFn(validator func(rune) bool) (accepted bool) {
	var next = lxr.Next()
	defer lxr.Backup()
	if next != EOF && !validator(next) {return false}
	accepted = true
	for next != EOF && validator(next) {
		next = lxr.Next()
	}
	return
}

func (lxr *lexer[T]) AcceptRune(rune rune) (accepted bool) {
	if accepted = lxr.Next() == rune; !accepted {
		lxr.Backup() 
	}
	return
}


func (lxr *lexer[T]) BufferedText() string  {
	return lxr.input[lxr.pos.Start:lxr.pos.End]
}

func (lxr *lexer[T]) RemainingText() string  {
	return lxr.input[lxr.pos.Start:]
}

func (lxr *lexer[T]) HasBufferedText() bool  {
	return lxr.pos.End > lxr.pos.Start
}

