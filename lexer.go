package e

import (
	"io"
	"text/scanner"
)

type Lexer struct {
	s scanner.Scanner
}

func NewLexer(src io.Reader) *Lexer {
	l := new(Lexer)
	l.s.Init(src)
	l.s.Error = func(s *scanner.Scanner, msg string) {
		panic(&SyntaxError{
			Msg:      msg,
			Position: Position{s.Position},
		})
	}
	l.s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.SkipComments | scanner.ScanComments
	return l
}

func (l *Lexer) Next() Token {
	s := &l.s
	tok := s.Scan()
	txt := s.TokenText()

	// symbols that require not peeking
	var ttype TType
	switch tok {
	case scanner.EOF:
		ttype = TEOF
	case scanner.Float:
		ttype = TNum
	case scanner.Ident:
		ttype = TIdent
	case scanner.Int:
		ttype = TNum
	case scanner.String:
		ttype = TString
	case '(':
		ttype = TLParen
	case ')':
		ttype = TRParen
	case '*':
		ttype = TMul
	case '+':
		ttype = TAdd
	case ',':
		ttype = TComma
	case '/':
		ttype = TDiv
	case '{':
		ttype = TLBrace
	case '}':
		ttype = TRBrace
	}
	if ttype != 0 {
		return Token{ttype, txt}
	}

	// symbols that require peeking
	peek := s.Peek()
	switch {
	case tok == '-' && peek == '>':
		ttype = TLambda
	case tok == '=' && peek == '=':
		ttype = TEquals
	case tok == '<' && peek == '=':
		ttype = TLe
	case tok == '>' && peek == '=':
		ttype = TGe
	}
	if ttype != 0 {
		s.Scan()
		return Token{ttype, txt + s.TokenText()}
	}

	// no peeked symbol was accepted
	switch {
	case tok == '-':
		ttype = TMinus
	case tok == '=':
		ttype = TAssign
	case tok == '<':
		ttype = TLt
	case tok == '>':
		ttype = TGt
	}
	if ttype != 0 {
		return Token{ttype, txt}
	}

	// no valid symbol was accepted
	panic(l.syntaxError("unexpected: " + scanner.TokenString(tok)))
}

func (l *Lexer) Position() Position {
	return Position{l.s.Position}
}

// returns a syntax error for the current position
func (l *Lexer) syntaxError(msg string) *SyntaxError {
	return &SyntaxError{
		Msg:      msg,
		Position: Position{l.s.Position},
	}
}
