// inspired by https://talks.golang.org/2011/lex.slide and https://golang.org/src/text/template/parse/lex.go

package main

import (
	"errors"
	"fmt"
	"strings"
)

// Lex turns a string into tokens.
func Lex(input string) ([]Token, error) {
	var tokens []Token
	for t := range lex(input) {
		if t.Type == tErr {
			return nil, errors.New(t.Value)
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

const (
	Digit      = "0123456789"
	Whitespace = " \t"
	EOF        = "\x00"
)

func lex(input string) chan Token {
	l := &lexer{
		input:  input,
		output: make(chan Token),
	}
	go l.run()
	return l.output
}

type lexer struct {
	input      string
	output     chan Token
	start, pos int
}

func (l *lexer) run() {
	for state := l.lexStart; state != nil; {
		state = state()
	}
	close(l.output)
}

type stateFn func() stateFn

func (l *lexer) lexStart() stateFn {
	p := l.peek()
	switch {
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

func (l *lexer) lexNum() stateFn {
	l.acceptN(Digit)
	l.emit(TNum)
	return l.lexStart
}

func (l *lexer) lexWhitespace() stateFn {
	l.acceptN(Whitespace)
	l.emitNone()
	return l.lexStart
}

func (l *lexer) lexEOF() stateFn {
	l.accept1(EOF)
	l.emit(TEOF)
	return nil
}

func (l *lexer) lexError(msg string) stateFn {
	l.emitError(fmt.Sprintf("pos %v: %v", l.pos, msg))
	return nil
}

//

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

//

func (l *lexer) emit(t Type) {
	// do not emit out-of-bounds
	stop := l.pos
	if stop > len(l.input) {
		stop = len(l.input)
	}

	l.output <- Token{t, l.input[l.start:stop]}
	l.start = l.pos
}

func (l *lexer) emitNone() {
	l.start = l.pos
}

func (l *lexer) emitError(msg string) {
	l.output <- Token{tErr, msg}
}

// is returns whether set contains x.
// E.g.:
// 	is('2', Digit) // true
// 	is('a', Digit) // false
func is(x byte, set string) bool {
	return strings.Contains(set, string(x))
}
