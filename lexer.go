// inspired by
// https://talks.golang.org/2011/lex.slide and
// https://golang.org/src/text/template/parse/lex.go

package main

import (
	"fmt"
	"strings"
)

// LexAll splits a string in tokens.
func LexAll(input string) (t []Token, e error) {
	defer func() {
		switch err := recover().(type) {
		default:
			panic(err) // resume
		case nil:
			// no error
		case SyntaxError:
			t = nil
			e = err
		}
	}()
	l := &lexer{input: input}
	var out []Token
	for {
		t := l.Next()
		switch t.TType {
		default:
			out = append(out, t)
		case TEOF, tErr:
			out = append(out, t)
			return out, nil
		}
	}
	return out, nil
}

// A Token represents a textual element like a word, number, ...
type Token struct {
	TType
	Value string
	//Pos   int
}

func (t Token) String() string {
	return t.TType.String() + "(" + t.Value + ")"
}

// TType enumerates all token types.
type TType int

// All possible token types
const (
	TEOF        TType = 0 // end-of-file
	TComma      TType = ','
	TLParen     TType = '('
	TRParen     TType = ')'
	TAdd        TType = '+'
	TDiv        TType = '/'
	TMinus      TType = '-'
	TMul        TType = '*'
	tErr        TType = 255         // error, internal use, filtered from output
	tWhitespace TType = tErr + iota // whitespace, internal use, filtered from output
	TIdent                          // identifier
	TNum                            // number
	TString                         // string
)

func (t TType) String() string {
	n, ok := ttypeName[t]
	if !ok {
		return fmt.Sprint("BAD_TYPE_", int(t))
	}
	return n
}

var ttypeName = map[TType]string{
	tErr:        "Err",
	tWhitespace: "whitespace",
	TComma:      ",",
	TEOF:        "EOF",
	TIdent:      "identifier",
	TLParen:     "(",
	TNum:        "number",
	TRParen:     ")",
	TString:     "string",
	TAdd:        "+",
	TDiv:        "/",
	TMinus:      "-",
	TMul:        "*",
}

// lex
//----------------------------------------------

const (
	EOF        = "\x00"
	Whitespace = " \t"
	Separator  = EOF + Whitespace

	NonZero = "123456789"
	Digit   = "0" + NonZero

	Lower    = "abcdefghijklmnopqrstuvwxyz"
	Upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alpha    = Lower + Upper
	AlphaNum = Alpha + Digit

	Delim = "(),+-*/"
	Quote = `"`
)

type lexer struct {
	input      string
	start, pos int
}

func (l *lexer) Next() Token {
	t := l.lexToken()
	if t.TType == tWhitespace {
		return l.Next()
	}
	return t
}

func (l *lexer) lexToken() Token {
	p := l.peek()
	switch {
	default:
		l.lexError("illegal character")
		panic("unreachable")
	case is(p, Alpha):
		return l.lexIdent()
	case is(p, Digit):
		return l.lexNum()
	case is(p, Whitespace):
		return l.lexWhitespace()
	case is(p, Delim):
		return l.lexDelim()
	case is(p, EOF):
		return l.lexEOF()
	case is(p, Quote):
		return l.lexString()
	}
}

func (l *lexer) lexNum() Token {
	l.acceptN(Digit)
	l.acceptN(AlphaNum) // accept trailing crap, atoi will catch this
	return l.emit(TNum)
}

func (l *lexer) lexIdent() Token {
	l.accept(Alpha)
	l.acceptN(AlphaNum)
	return l.emit(TIdent)
}

func (l *lexer) lexWhitespace() Token {
	l.acceptN(Whitespace)
	return l.emit(tWhitespace)
}

func (l *lexer) lexDelim() Token {
	l.accept(Delim)
	t := l.emit(0)
	t.TType = TType(t.Value[0])
	return t
}

func (l *lexer) lexString() Token {
	l.accept(Quote)
	p := l.peek()
	for p != '"' && p != 0 {
		l.acceptAny()
		p = l.peek()
	}

	// TODO: expect
	if l.peek() != '"' {
		l.lexError("unterminated string")
	}
	l.accept(Quote)

	return l.emit(TString)
}

func (l *lexer) lexEOF() Token {
	l.accept(EOF)
	return l.emit(TEOF)
}

func (l *lexer) lexError(msg string) {
	panic(SyntaxError{fmt.Sprint("pos %v: %v", l.pos, msg)})
}

//----------------------------------------------

func (l *lexer) peek() byte {
	if l.pos >= len(l.input) {
		return bEOF
	}
	return l.input[l.pos]
}

const bEOF = 0

func (l *lexer) acceptN(set string) {
	for l.accept(set) {
	}
}

func (l *lexer) accept(set string) bool {
	if is(l.peek(), set) {
		l.pos++
		return true
	}
	return false
}

func (l *lexer) acceptAny() {
	l.pos++
}

//----------------------------------------------

// emit returns a token for the current position,
// and advances the position.
func (l *lexer) emit(t TType) Token {
	// do not emit out-of-bounds
	stop := l.pos
	if stop > len(l.input) {
		stop = len(l.input)
	}
	tok := Token{t, l.input[l.start:stop]}
	l.start = l.pos
	return tok
}

//----------------------------------------------

// is returns whether set contains x. E.g.:
// 	is('2', Digit) // true
// 	is('a', Digit) // false
func is(x byte, set string) bool {
	return strings.Contains(set, string(x))
}