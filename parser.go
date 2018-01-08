package se

import (
	"fmt"
	"io"
	"strconv"
)

func Parse(src io.Reader) (Node, error) {
	return NewParser(src).Parse()
}

type Parser struct {
	lex  Lexer
	next [readAhead]Token
}

const readAhead = 4

func NewParser(src io.Reader) *Parser {
	return &Parser{lex: *NewLexer(src)}
}

func (p *Parser) Parse() (_ Node, e error) {
	return withCatch(func() Node {
		p.init()
		program := p.PExpr()
		p.Expect(TEOF)
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
			case *SyntaxError:
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
	if p.HasPeek(TLParen, TRParen) ||
		p.HasPeek(TLParen, TIdent, TComma) ||
		p.HasPeek(TIdent, TLambda) ||
		p.HasPeek(TLParen, TIdent, TRParen, TLambda) {
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
	if p.HasPeek(TIdent) {
		args = []*Ident{p.PIdent()}
	} else {
		args = p.PIdentList()
	}

	p.Expect(TLambda)

	body := p.PExpr()
	return &Lambda{Args: args, Body: body}
}

// identlist:
//  | ()
//  | (ident,...)
func (p *Parser) PIdentList() []*Ident {
	p.Expect(TLParen)
	var l []*Ident

	//()
	if p.Accept(TRParen) {
		return l
	}

	//(ident,...)
	l = append(l, p.PIdent())
	for p.Accept(TComma) {
		l = append(l, p.PIdent())
	}
	p.Expect(TRParen)
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

func opFunc(t TType) string {
	if f, ok := opfunc[t]; ok {
		return f
	}
	panic(fmt.Sprintf("bug: bad operator: %v", t))
}

var opfunc = map[TType]string{
	TAdd: "add",
	TMul: "mul",
}

var isUnary = map[TType]bool{
	TAdd:   true,
	TMinus: true,
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
	if p.Accept(TMinus) {
		return &Call{&Ident{Name: "neg"}, []Node{p.POperand()}}
	}

	// num, ident, parenexpr
	var expr Node
	switch p.PeekTT() {
	case TNum:
		expr = p.PNum()
	case TIdent:
		expr = p.PIdent()
	case TLParen:
		expr = p.PParenExpr()
	default:
		panic(p.Unexpected(p.Next()))
	}

	// operand *(list): function call
	for p.PeekTT() == TLParen {
		args := p.PArgList()
		expr = &Call{expr, args}
	}

	return expr
}

// parse a number.
func (p *Parser) PNum() Node {
	tok := p.Expect(TNum)
	v, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		panic(p.SyntaxError(err.Error()))
	}
	return &Num{v}
}

// parse an identifier
func (p *Parser) PIdent() *Ident {
	tok := p.Expect(TIdent)
	return &Ident{Name: tok.Value}
}

// parse a parenthesized argument list:
//  arglist:
//   | ()
//   | ( expr1, expr1, ... )
func (p *Parser) PArgList() []Node {
	p.Expect(TLParen)

	// ()
	if p.Accept(TRParen) {
		return []Node{}
	}

	// ( expr1, expr1, ... )
	list := []Node{p.PExpr1()}
	for p.Accept(TComma) {
		list = append(list, p.PExpr1())
	}
	p.Expect(TRParen)
	return list
}

func (p *Parser) PParenExpr() Node {
	p.Expect(TLParen)
	//if p.Accept(TRParen) {
	//	return List{}
	//}
	expr := p.PExpr()
	p.Expect(TRParen)
	return expr
}

// ------------------------------------------

var precedence = map[TType]int{
	TAdd:   1,
	TMinus: 1,
	TMul:   2,
	TDiv:   2,
}

// Peek returns the next token in the stream without advancing
func (p *Parser) Peek() Token {
	return p.next[0]
}

func (p *Parser) HasPeek(want ...TType) bool {
	for i, w := range want {
		if p.next[i].TType != w {
			return false
		}
	}
	return true
}

func (p *Parser) PeekTT() TType {
	return p.Peek().TType
}

// Next returns the next token in the stream and advances
func (p *Parser) Next() Token {
	curr := p.next[0]

	for i := 0; i < readAhead-1; i++ {
		p.next[i] = p.next[i+1]
	}
	p.next[readAhead-1] = p.lex.Next()
	return curr
}

// if the peeked token is of type t, consume the token and return true.
func (p *Parser) Accept(t TType) bool {
	if p.Peek().TType == t {
		p.Next()
		return true
	}
	return false
}

// consume the next token and throw an error if it is not of the expected type.
func (p *Parser) Expect(t TType) Token {
	if n := p.Next(); n.TType != t {
		panic(p.SyntaxError(fmt.Sprintf("unexpected '%v', expected '%v'", n, t)))
	} else {
		return n
	}
}

// construct a syntax error for unexpected token at current position.
func (p *Parser) Unexpected(t Token) *SyntaxError {
	return p.SyntaxError(fmt.Sprintf("unexpected '%v'", t))
}

// construct a syntax error at current position.
func (p *Parser) SyntaxError(msg string) *SyntaxError {
	return &SyntaxError{Msg: msg, Position: p.nextPos()}
}

func (p *Parser) nextPos() Position {
	// TODO: return p.Peek().Pos
	return Position{}
}

func (p *Parser) init() {
	if p.next[0] != (Token{}) {
		panic("parser: init called twice")
	}
	for i := 0; i < readAhead; i++ {
		p.Next()
	}
}
