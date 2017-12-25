package main

import (
	"fmt"
	"strconv"
)

func Parse(in string) (ex Expr, err error) {
	// catch errors
	defer func() {
		if e := recover(); e != nil {
			if parseErr, ok := e.(SyntaxError); ok {
				err = parseErr
			} else {
				panic(e) // re-throw
			}
		}
	}()

	// lex
	tokens, err := Lex(in)
	if err != nil {
		return nil, err
	}

	// parse
	p := &parser{input: tokens}
	return p.parseProgram(), nil
}

type parser struct {
	input []Token
	pos   int
}

func (p *parser) parseProgram() Expr {
	e := p.parseExpr()
	p.expect(TEOF)
	return e
}

// parseExpr parses an expression, like:
// 	(1+2)*3
// 	f(x,y)
// 	f()()
func (p *parser) parseExpr() Expr {
	var expr Expr

	// non-call expression
	t := p.peek()
	switch t.TType {
	default:
		p.errorf("unexpected %v", t.TType)
		return nil // unreachable
	case TNum:
		expr = p.parseNum()
	case TLParen:
		expr = p.parseParenExpr()
	case TIdent:
		expr = p.parseIdent()
	}

	// call expression
	for p.peek().TType == TLParen {
		args := p.parseArgs()
		expr = &Call{expr, args}
	}
	return expr
}

// parseArgs parses an argument list, like:
// 	()
// 	(x, y)
func (p *parser) parseArgs() []Expr {
	var args []Expr

	p.expect(TLParen)

	for {
		if p.peek().TType == TRParen {
			p.accept(TRParen)
			return args
		}
		args = append(args, p.parseExpr())
		if p.peek().TType != TRParen {
			p.expect(TComma)
		}
	}
	panic("unreachable")
}

// parseIdent parses an identifier, like:
// 	foo
func (p *parser) parseIdent() Expr {
	tok := p.expect(TIdent)
	return &Ident{tok.Value}
}

// parseNum parses a number, like:
// 	42
func (p *parser) parseNum() Expr {
	tok := p.expect(TNum)
	v, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		p.errorf("%v", err)
	}
	return &Num{v}
}

// parseParenExpr parses a parenthesised expression, like:
// 	(1+x)
func (p *parser) parseParenExpr() Expr {
	p.expect(TLParen)
	e := p.parseExpr()
	p.expect(TRParen)
	return e
}

//----------------------------------------------

func (p *parser) peek() Token {
	if p.pos >= len(p.input) {
		panic("BUG: beyond input")
	}
	return p.input[p.pos]
}

func (p *parser) expect(typ TType) Token {
	tok := p.peek()
	if tok.TType != typ {
		p.errorf("unexpected %v, expected %v", tok.TType, typ)
	}
	p.advance()
	return tok
}

func (p *parser) accept(typ TType) Token {
	tok := p.peek()
	p.advance()
	return tok
}

func (p *parser) advance() {
	p.pos++
}

func (p *parser) errorf(format string, x ...interface{}) {
	panic(SyntaxError{fmt.Sprintf("pos %v: %v", p.pos, fmt.Sprintf(format, x...))})
}
