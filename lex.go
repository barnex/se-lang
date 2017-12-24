// inspired by
// https://talks.golang.org/2011/lex.slide and
// https://golang.org/src/text/template/parse/lex.go

package main

import (
	"errors"
	"fmt"
	"strings"
)

// Lex splits a string in tokens.
func Lex(input string) ([]Token, error) {
	out := (&lexer{input: input}).lexAll()
	// if we have an error token, return error explicitly
	if last := out[len(out)-1]; last.TType == tErr {
		return nil, errors.New(last.Value)
	}
	return out, nil
}

// token
//----------------------------------------------

// A Token represents a textual element like a word, number, ...
type Token struct {
	TType
	Value string
}

// TType enumerates all token types.
type TType int

// All possible token types
const (
	tErr        TType = iota // error, internal use, filtered from output
	tWhitespace              // whitespace, internal use, filtered from output
	TEOF                     // end-of-file
	TIdent                   // identifier
	TNum                     // number
)

func (t Token) String() string {
	return t.TType.String() + "(" + t.Value + ")"
}

func (t TType) String() string {
	n, ok := tokenName[t]
	if !ok {
		return fmt.Sprint("BAD_TYPE_", int(t))
	}
	return n
}

var tokenName = map[TType]string{
	tErr:        "Err",
	tWhitespace: "Whitespace",
	TEOF:        "EOF",
	TIdent:      "Ident",
	TNum:        "Num",
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
)

type lexer struct {
	input      string
	start, pos int
}

func (l *lexer) lexAll() []Token {
	var out []Token
	for {
		t := l.lexToken()
		switch t.TType {
		default:
			out = append(out, t)
		case tWhitespace: // ignore
		case TEOF, tErr:
			out = append(out, t)
			return out
		}
	}
	return out
}

func (l *lexer) lexToken() Token {
	p := l.peek()
	switch {
	case is(p, Alpha):
		return l.lexIdent()
	case is(p, Digit):
		return l.lexNum()
	case is(p, Whitespace):
		return l.lexWhitespace()
	case is(p, EOF):
		return l.lexEOF()
	default:
		return l.lexError("illegal character")
	}
}

func (l *lexer) lexNum() Token {
	l.acceptN(Digit)
	l.acceptN(AlphaNum) // accept trailing crap, atoi will catch this
	return l.emit(TNum)
}

func (l *lexer) lexIdent() Token {
	l.accept1(Alpha)
	l.acceptN(AlphaNum)
	return l.emit(TIdent)
}

func (l *lexer) lexWhitespace() Token {
	l.acceptN(Whitespace)
	return l.emit(tWhitespace)
}

func (l *lexer) lexEOF() Token {
	l.accept1(EOF)
	return l.emit(TEOF)
}

func (l *lexer) lexError(msg string) Token {
	tok := l.emit(tErr)
	tok.Value = fmt.Sprintf("pos %v: %q: %v", l.pos, tok.Value, msg)
	return tok
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
	for l.accept1(set) {
	}
}

func (l *lexer) accept1(set string) bool {
	if is(l.peek(), set) {
		l.accept()
		return true
	}
	return false
}

func (l *lexer) accept() {
	l.pos++
}

//----------------------------------------------

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

// is returns whether set contains x.
// E.g.:
// 	is('2', Digit) // true
// 	is('a', Digit) // false
func is(x byte, set string) bool {
	return strings.Contains(set, string(x))
}
