package se

import (
	"fmt"
	"io"
	"strconv"

	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/lex"
)

func Parse(src io.Reader) (Node, error) {
	return NewParser(src).Parse()
}

type Parser struct {
	lex  lex.Lexer
	next [readAhead]lex.Token
}

const readAhead = 4

func NewParser(src io.Reader) *Parser {
	return &Parser{lex: *lex.NewLexer(src)}
}

func (p *Parser) Parse() (_ Node, e error) {
	return withCatch(func() Node {
		p.init()
		program := p.PExpr()
		p.Expect(lex.TEOF)
		return program
	})
}

// debug: panic on parse error
const panicOnErr = false

func withCatch(f func() Node) (_ Node, e error) {
	// catch syntax errors
	if !panicOnErr {
		defer func() {
			switch err := recover().(type) {
			default:
				panic(err) // resume
			case nil:
				// no error
			case se.Error:
				e = err
			}
		}()
	}
	return f(), nil
}

// --------

// expr:
// 	| expr1
//  | lambda
func (p *Parser) PExpr() Node {
	// peek for lambda: "()" or "(ident," or "ident->" or "(ident)->"
	if p.HasPeek(lex.TLParen, lex.TRParen) ||
		p.HasPeek(lex.TLParen, lex.TIdent, lex.TComma) ||
		p.HasPeek(lex.TIdent, lex.TLambda) ||
		p.HasPeek(lex.TLParen, lex.TIdent, lex.TRParen, lex.TLambda) {
		return p.PLambda()
	} else {
		return p.PExpr1()
	}
}

// lambda:
//  | ident -> expr1
//  | () -> expr1
//  | (ident) -> expr1
//  | (ident,...) -> expr1
func (p *Parser) PLambda() Node {
	var args []*Ident

	// ident -> expr
	if p.HasPeek(lex.TIdent) {
		args = []*Ident{p.PIdent()}
	} else {
		args = p.PIdentList()
	}

	p.Expect(lex.TLambda)

	body := p.PExpr()
	return &Lambda{Args: args, Body: body}
}

// identlist:
//  | ()
//  | (ident,...)
func (p *Parser) PIdentList() []*Ident {
	p.Expect(lex.TLParen)
	var l []*Ident

	//()
	if p.Accept(lex.TRParen) {
		return l
	}

	//(ident,...)
	l = append(l, p.PIdent())
	for p.Accept(lex.TComma) {
		l = append(l, p.PIdent())
	}
	p.Expect(lex.TRParen)
	return l
}

// PExpr parses an expression not containing lambdas
//  expr:
//   | operand                      // expression without infix operators
//   | operand operator expr1       // binary operator
func (p *Parser) PExpr1() Node {
	return p.PBinary(1)
}

// parse an expression, or binary expression as long as operator precedence is at least prec1.
// inspired by https://github.com/adonovan/gopl.io/blob/master/ch7/eval/parse.go
func (p *Parser) PBinary(prec1 int) Node {
	lhs := p.POperand()
	for prec := precedence[p.Peek().TType]; prec >= prec1; prec-- {
		for precedence[p.Peek().TType] == prec {
			op := p.Next()
			rhs := p.PBinary(prec + 1)
			lhs = &Call{&Ident{Name: opFunc(op.TType)}, []Node{lhs, rhs}}
		}
	}
	return lhs
}

func opFunc(t lex.TType) string {
	if f, ok := opfunc[t]; ok {
		return f
	}
	panic(fmt.Sprintf("bug: bad operator: %v", t))
}

var opfunc = map[lex.TType]string{
	lex.TAdd: "add",
	lex.TMul: "mul",
}

var isUnary = map[lex.TType]bool{
	lex.TAdd:   true,
	lex.TMinus: true,
}

// parse an operand:
// operand:
//  | - operand
//  | num
//  | ident
//  | parenexpr
//  | operand *(list)
func (p *Parser) POperand() Node {

	// - operand
	if p.Accept(lex.TMinus) {
		return &Call{&Ident{Name: "neg"}, []Node{p.POperand()}}
	}

	// num, ident, parenexpr
	var expr Node
	switch p.PeekTT() {
	case lex.TNum:
		expr = p.PNum()
	case lex.TIdent:
		expr = p.PIdent()
	case lex.TLParen:
		expr = p.PParenExpr()
	default:
		panic(p.Unexpected(p.Next()))
	}

	// operand *(list): function call
	for p.PeekTT() == lex.TLParen {
		args := p.PArgList()
		expr = &Call{expr, args}
	}

	return expr
}

// parse a number.
func (p *Parser) PNum() Node {
	tok := p.Expect(lex.TNum)
	v, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		panic(p.SyntaxError(err.Error()))
	}
	return &Num{v}
}

// parse an identifier
func (p *Parser) PIdent() *Ident {
	tok := p.Expect(lex.TIdent)
	return &Ident{Name: tok.Value}
}

// parse a parenthesized argument list:
//  arglist:
//   | ()
//   | ( expr1, expr1, ... )
func (p *Parser) PArgList() []Node {
	p.Expect(lex.TLParen)

	// ()
	if p.Accept(lex.TRParen) {
		return []Node{}
	}

	// ( expr1, expr1, ... )
	list := []Node{p.PExpr1()}
	for p.Accept(lex.TComma) {
		list = append(list, p.PExpr1())
	}
	p.Expect(lex.TRParen)
	return list
}

func (p *Parser) PParenExpr() Node {
	p.Expect(lex.TLParen)
	//if p.Accept(TRParen) {
	//	return List{}
	//}
	expr := p.PExpr()
	p.Expect(lex.TRParen)
	return expr
}

// ------------------------------------------

var precedence = map[lex.TType]int{
	lex.TAdd:   1,
	lex.TMinus: 1,
	lex.TMul:   2,
	lex.TDiv:   2,
}

// Peek returns the next token in the stream without advancing
func (p *Parser) Peek() lex.Token {
	return p.next[0]
}

func (p *Parser) HasPeek(want ...lex.TType) bool {
	for i, w := range want {
		if p.next[i].TType != w {
			return false
		}
	}
	return true
}

func (p *Parser) PeekTT() lex.TType {
	return p.Peek().TType
}

// Next returns the next token in the stream and advances
func (p *Parser) Next() lex.Token {
	curr := p.next[0]

	for i := 0; i < readAhead-1; i++ {
		p.next[i] = p.next[i+1]
	}
	p.next[readAhead-1] = p.lex.Next()
	return curr
}

// if the peeked token is of type t, consume the token and return true.
func (p *Parser) Accept(t lex.TType) bool {
	if p.Peek().TType == t {
		p.Next()
		return true
	}
	return false
}

// consume the next token and throw an error if it is not of the expected type.
func (p *Parser) Expect(t lex.TType) lex.Token {
	if n := p.Next(); n.TType != t {
		panic(p.SyntaxError(fmt.Sprintf("unexpected '%v', expected '%v'", n, t)))
	} else {
		return n
	}
}

// construct a syntax error for unexpected token at current position.
func (p *Parser) Unexpected(t lex.Token) se.Error {
	return p.SyntaxError(fmt.Sprintf("unexpected '%v'", t))
}

// construct a syntax error at current position.
func (p *Parser) SyntaxError(msg string) se.Error {
	return se.Errorf("line %v: %v", p.nextPos(), msg)
}

func (p *Parser) nextPos() se.Position {
	// TODO: return p.Peek().Pos
	return se.Position{}
}

func (p *Parser) init() {
	if p.next[0] != (lex.Token{}) {
		panic("parser: init called twice")
	}
	for i := 0; i < readAhead; i++ {
		p.Next()
	}
}
