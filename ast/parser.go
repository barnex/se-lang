/*
	Package ast implements the language's parser.
	I.e. parses source text into an Abstract Syntax Tree.
*/
package ast

import (
	"fmt"
	"io"
	"strconv"

	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/lex"
)

// Parse reads source text from src
// and returns an Abstract Syntax Tree.
func Parse(src io.Reader) (_ Node, e error) {
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
	p := parser{lex: *lex.NewLexer(src)}
	return p.parse(), nil
}

// The syntax is LL(4),
// mainly due to the lambda syntactic sugar
// 	(x,y)->...
// 	x->...
const readAhead = 4

type parser struct {
	lex  lex.Lexer
	next [readAhead]lex.Token
}

func (p *parser) parse() Node {
	p.init()
	program := p.parseExpr()
	p.Expect(lex.TEOF)
	return program
}

// expr:
// 	| expr1
//  | lambda
//  | block
func (p *parser) parseExpr() Node {
	// peek for lambda: "()" or "(ident," or "ident->" or "(ident)->"
	if p.HasPeek(lex.TLParen, lex.TRParen) ||
		p.HasPeek(lex.TLParen, lex.TIdent, lex.TComma) ||
		p.HasPeek(lex.TIdent, lex.TLambda) ||
		p.HasPeek(lex.TLParen, lex.TIdent, lex.TRParen, lex.TLambda) {
		return p.parseLambda()
	}

	if p.HasPeek(lex.TLBrace) {
		return p.parseBlock()
	}

	return p.parseExpr1()
}

// block:
//  | { stmt; ... }
func (p *parser) parseBlock() Node {
	p.Expect(lex.TLBrace)
	stmt := []Node{p.parseStmt()}

	for p.Accept(lex.TSemicol) {
		stmt = append(stmt, p.parseStmt())
	}
	p.Expect(lex.TRBrace)
	return &Block{stmt}
}

// stmt:
// 	| expr
//  | assign
func (p *parser) parseStmt() Node {
	if p.HasPeek(lex.TIdent, lex.TAssign) {
		return p.parseAssign()
	} else {
		return p.parseExpr()
	}
}

// assign:
//  | ident = expr
func (p *parser) parseAssign() Node {
	lhs := p.parseIdent()
	p.Expect(lex.TAssign)
	rhs := p.parseExpr()
	return &Assign{lhs, rhs}
}

// lambda:
//  | ident -> expr1
//  | () -> expr1
//  | (ident) -> expr1
//  | (ident,...) -> expr1
func (p *parser) parseLambda() Node {
	var args []*Ident

	// ident -> expr
	if p.HasPeek(lex.TIdent) {
		args = []*Ident{p.parseIdent()}
	} else {
		args = p.parseIdentList()
	}

	p.Expect(lex.TLambda)

	body := p.parseExpr()
	return &Lambda{Args: args, Body: body}
}

// identlist:
//  | ()
//  | (ident,...)
func (p *parser) parseIdentList() []*Ident {
	p.Expect(lex.TLParen)
	var l []*Ident

	//()
	if p.Accept(lex.TRParen) {
		return l
	}

	//(ident,...)
	l = append(l, p.parseIdent())
	for p.Accept(lex.TComma) {
		l = append(l, p.parseIdent())
	}
	p.Expect(lex.TRParen)
	return l
}

// parses an expression not containing lambdas
//  expr1:
//   | operand                      // expression without infix operators
//   | operand operator expr1       // binary operator
func (p *parser) parseExpr1() Node {
	return p.parseBinaryExpr(1)
}

// parse an expression, or binary expression as long as operator precedence is at least prec1.
// inspired by https://github.com/adonovan/gopl.io/blob/master/ch7/eval/parse.go
func (p *parser) parseBinaryExpr(prec1 int) Node {
	lhs := p.parseOperand()
	for prec := precedence[p.Peek().TType]; prec >= prec1; prec-- {
		for precedence[p.Peek().TType] == prec {
			op := p.Next()
			rhs := p.parseBinaryExpr(prec + 1)
			lhs = &Call{&Ident{Name: opFunc(op.TType)}, []Node{lhs, rhs}}
		}
	}
	return lhs
}

// parse an operand:
// operand:
//  | - operand
//  | num
//  | ident
//  | parenexpr
//  | operand *(list)
func (p *parser) parseOperand() Node {

	// - operand
	if p.Accept(lex.TMinus) {
		return &Call{&Ident{Name: "neg"}, []Node{p.parseOperand()}}
	}

	// !operand
	if p.Accept(lex.TNot) {
		return &Call{&Ident{Name: "not"}, []Node{p.parseOperand()}}
	}

	// num, ident, parenexpr
	var expr Node
	switch p.PeekTT() {
	case lex.TNum:
		expr = p.parseNum()
	case lex.TIdent:
		expr = p.parseIdent()
	case lex.TLParen:
		expr = p.parseParenExpr()
	default:
		panic(p.Unexpected(p.Next()))
	}

	// operand *(list): function call
	for p.PeekTT() == lex.TLParen {
		args := p.parseArgList()
		expr = &Call{expr, args}
	}

	return expr
}

// parse a number.
func (p *parser) parseNum() Node {
	tok := p.Expect(lex.TNum)
	v, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		panic(p.SyntaxError(err.Error()))
	}
	return &Num{v}
}

// parse an identifier
func (p *parser) parseIdent() *Ident {
	tok := p.Expect(lex.TIdent)
	return &Ident{Name: tok.Value}
}

// parse a parenthesized argument list:
//  arglist:
//   | ()
//   | ( expr1, expr1, ... )
func (p *parser) parseArgList() []Node {
	p.Expect(lex.TLParen)

	// ()
	if p.Accept(lex.TRParen) {
		return []Node{}
	}

	// ( expr1, expr1, ... )
	list := []Node{p.parseExpr1()}
	for p.Accept(lex.TComma) {
		list = append(list, p.parseExpr1())
	}
	p.Expect(lex.TRParen)
	return list
}

func (p *parser) parseParenExpr() Node {
	p.Expect(lex.TLParen)
	expr := p.parseExpr()
	p.Expect(lex.TRParen)
	return expr
}

// ------------------------------------------

var precedence = map[lex.TType]int{
	lex.TMul: 5,
	lex.TDiv: 5,
	lex.TMod: 5,

	lex.TAdd:   4,
	lex.TMinus: 4,

	lex.TEq:  3,
	lex.TGe:  3,
	lex.TGt:  3,
	lex.TLe:  3,
	lex.TLt:  3,
	lex.TNEq: 3,

	lex.TAnd: 2,

	lex.TOr: 1,
}

func opFunc(t lex.TType) string {
	if f, ok := opStr[t]; ok {
		return f
	}
	panic(fmt.Sprintf("bug: bad operator: %v", t))
}

var opStr = map[lex.TType]string{
	lex.TAdd:   "add",
	lex.TAnd:   "and",
	lex.TEq:    "eq",
	lex.TGe:    "ge",
	lex.TGt:    "gt",
	lex.TLe:    "le",
	lex.TLt:    "lt",
	lex.TMinus: "sub",
	lex.TMod:   "mod",
	lex.TMul:   "mul",
	lex.TNEq:   "neq",
	lex.TOr:    "or",
}

var isUnary = map[lex.TType]bool{
	lex.TAdd:   true,
	lex.TMinus: true,
	lex.TNot:   true,
}

// ------------------------------------------

// Peek returns the next token in the stream without advancing
func (p *parser) Peek() lex.Token {
	return p.next[0]
}

func (p *parser) HasPeek(want ...lex.TType) bool {
	for i, w := range want {
		if p.next[i].TType != w {
			return false
		}
	}
	return true
}

func (p *parser) PeekTT() lex.TType {
	return p.Peek().TType
}

// Next returns the next token in the stream and advances
func (p *parser) Next() lex.Token {
	curr := p.next[0]

	for i := 0; i < readAhead-1; i++ {
		p.next[i] = p.next[i+1]
	}
	p.next[readAhead-1] = p.lex.Next()
	return curr
}

// if the peeked token is of type t, consume the token and return true.
func (p *parser) Accept(t lex.TType) bool {
	if p.Peek().TType == t {
		p.Next()
		return true
	}
	return false
}

// consume the next token and throw an error if it is not of the expected type.
func (p *parser) Expect(t lex.TType) lex.Token {
	if n := p.Next(); n.TType != t {
		panic(p.SyntaxError(fmt.Sprintf("unexpected '%v', expected '%v'", n, t)))
	} else {
		return n
	}
}

// construct a syntax error for unexpected token at current position.
func (p *parser) Unexpected(t lex.Token) se.Error {
	return p.SyntaxError(fmt.Sprintf("unexpected '%v'", t))
}

// construct a syntax error at current position.
func (p *parser) SyntaxError(msg string) se.Error {
	return se.Errorf("line %v: %v", p.nextPos(), msg)
}

func (p *parser) nextPos() se.Position {
	// TODO: return p.Peek().Pos
	return se.Position{}
}

func (p *parser) init() {
	if p.next[0] != (lex.Token{}) {
		panic("parser: init called twice")
	}
	for i := 0; i < readAhead; i++ {
		p.Next()
	}
}
