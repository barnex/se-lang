package e

import (
	"fmt"
	"io"
	"strconv"
)

func Parse(src io.Reader) (Node, error) {
	return NewParser(src).Parse()
}

type Parser struct {
	lex     Lexer
	next    Token
	nextPos Position
}

func NewParser(src io.Reader) *Parser {
	return &Parser{lex: *NewLexer(src)}
}

// debug: panic on parse error
const panicOnErr = false

func (p *Parser) Parse() (_ Node, e error) {
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

	p.init()
	program := p.PExpr()
	p.Expect(TEOF)
	return program, nil
}

// --------

// expr:
// 	| expr1, expr1,...
func (p *Parser) PExpr() Node {
	first := p.PExpr1()

	if p.PeekTT() != TComma {
		return first
	}

	expr := List{first}
	for p.Accept(TComma) {
		expr = append(expr, p.PExpr1())
	}
	return expr
}

// PExpr parses a single expression: not containing comma's
//  expr:
//   | operand                      // expression without infix operators
//   | operand operator expr1       // binary operator
func (p *Parser) PExpr1() Node {
	expr := p.PBinary(1)

	// lambda: assure argument list are identifiers
	if call, ok := expr.(List); ok {
		if id, ok := call.Car().(*Ident); ok {
			if id.Name == "lambda" {
				call[1] = p.toIdentList(call[1])
			}
		}
	}

	return expr
}

func (p *Parser) toIdentList(args Node) Node {
	var list List
	if l, ok := args.(List); ok {
		list = l
	} else {
		list = List{args}
	}
	for _, a := range list {
		if _, ok := a.(*Ident); !ok {
			panic(SyntaxErrorf("lambda: arguments must be identifiers, have %v", ToString(a)))
		}
	}
	return list
}

// parse an expression, or binary expression as long as operator precedence is at least prec1.
// inspired by https://github.com/adonovan/gopl.io/blob/master/ch7/eval/parse.go
func (p *Parser) PBinary(prec1 int) Node {
	lhs := p.POperand()
	for prec := precedence[p.Peek().TType]; prec >= prec1; prec-- {
		for precedence[p.Peek().TType] == prec {
			op := p.Next()
			rhs := p.PBinary(prec + 1)
			lhs = MakeList(&Ident{opFunc(op.TType)}, lhs, rhs)
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
	TLambda: "lambda",
	TAdd:    "add",
	TMul:    "mul",
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
		return List{&Ident{"neg"}, p.POperand()}
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

	// operand *(list)
	for p.PeekTT() == TLParen {
		args := p.PArgList()
		expr = MakeList(expr, args...)
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
func (p *Parser) PIdent() Node {
	tok := p.Expect(TIdent)
	return &Ident{tok.Value}
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
	if p.Accept(TRParen) {
		return List{}
	}
	expr := p.PExpr()
	p.Expect(TRParen)
	return expr
}

// ------------------------------------------

var precedence = map[TType]int{
	TLambda: 1,
	TAdd:    2,
	TMinus:  2,
	TMul:    3,
	TDiv:    3,
}

// Peek returns the next token in the stream without advancing
func (p *Parser) Peek() Token {
	return p.next
}

func (p *Parser) PeekTT() TType {
	return p.next.TType
}

// Next returns the next token in the stream and advances
func (p *Parser) Next() Token {
	curr := p.next
	p.nextPos = p.lex.Position()
	p.next = p.lex.Next()
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
	return &SyntaxError{Msg: msg, Position: p.nextPos}
}

func (p *Parser) init() {
	if p.next != (Token{}) {
		panic("parser: init called twice")
	}
	p.Next()
}
