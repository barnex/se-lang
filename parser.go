package main

import (
	"fmt"
	"strconv"
)

type Expr interface{}

type Num struct {
	Value float64
}

type Call struct {
	Op   Expr
	Args []Expr
}

func Parse(in string) (ex Expr, err error) {
	// catch errors
	defer func() {
		if e := recover(); e != nil {
			if parseErr, ok := e.(parseError); ok {
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

func (p *parser) parseExpr() Expr {
	t := p.peek()
	switch t.TType {
	default:
		p.errorf("unexpected %v", t.TType)
		return nil // unreachable
	case TNum:
		return p.parseNum()
	case TLParen:
		return p.parseParenExpr()
	}
}

func (p *parser) parseNum() Expr {
	tok := p.expect(TNum)
	v, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		p.errorf("%v", err)
	}
	return Num{v}
}

func (p *parser) parseParenExpr() Expr {
	p.expect(TLParen)
	e := p.parseExpr()
	p.expect(TRParen)
	return e
}

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

func (p *parser) advance() {
	p.pos++
}

func (p *parser) errorf(format string, x ...interface{}) {
	panic(parseError(fmt.Sprintf("pos %v: %v", p.pos, fmt.Sprintf(format, x...))))
}

type parseError string

func (e parseError) Error() string { return string(e) }

//func (p *parser) accept(typ TType) bool {
//
//}

//----------------------------------------------
