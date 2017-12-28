package e

import (
	"fmt"
	"io"
	"strconv"
)

func Parse(src io.Reader) (Expr, error) {
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

func (p *Parser) Parse() (_ Expr, e error) {
	// catch syntax errors
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

	p.init()
	program := p.parseExpr()
	p.Expect(TEOF)
	return program, nil
}

// parse an expression
func (p *Parser) parseExpr() Expr {
	return p.parseBinary(1)
}

var precedence = map[TType]int{
	TLambda: 1,
	TAdd:    2,
	TMinus:  2,
	TMul:    3,
	TDiv:    3,
}

// parse an expression, or binary expression as long as operator precedence is at least prec1.
// inspired by https://github.com/adonovan/gopl.io/blob/master/ch7/eval/parse.go
func (p *Parser) parseBinary(prec1 int) Expr {
	lhs := p.parsePrimary()
	for prec := precedence[p.Peek().TType]; prec >= prec1; prec-- {
		for precedence[p.Peek().TType] == prec {
			op := p.Next()
			rhs := p.parseBinary(prec + 1)
			lhs = &Call{
				Func: &Ident{op.Value},
				Args: []Expr{lhs, rhs},
			}
		}
	}
	return lhs
}

// parse an expression that does not contain binary operators.
func (p *Parser) parsePrimary() Expr {
	var expr Expr

	// non-call expression
	t := p.Peek()
	switch t.TType {
	default:
		panic(p.Unexpected(p.Next()))
	case TNum:
		expr = p.parseNum()
	case TLParen:
		expr = p.parseParenExpr()
	case TIdent:
		expr = p.parseIdent()
	}

	// call expression
	for p.Accept(TLParen) {
		args := p.parseArgs()
		p.Expect(TRParen)
		expr = &Call{expr, args}
	}

	return expr
}

// parse a function argument list
func (p *Parser) parseArgs() []Expr {
	var args []Expr

	for {
		if p.Peek().TType == TRParen {
			return args
		}
		args = append(args, p.parseExpr())
		if p.Peek().TType != TRParen {
			p.Expect(TComma)
		}
	}
	panic("unreachable")
}

// parse a parenthesized expression, stripping the outermost parens.
func (p *Parser) parseParenExpr() Expr {
	p.Expect(TLParen)
	e := p.parseExpr()
	p.Expect(TRParen)
	return e
}

// parse a number
func (p *Parser) parseNum() Expr {
	tok := p.Expect(TNum)
	v, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		panic(p.SyntaxError(err.Error()))
	}
	return &Num{v}
}

// parse an identifier
func (p *Parser) parseIdent() Expr {
	tok := p.Expect(TIdent)
	return &Ident{tok.Value}
}

// Peek returns the next token in the stream without advancing
func (p *Parser) Peek() Token {
	return p.next
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
		panic(p.SyntaxError(fmt.Sprintf("unexpected %v, expected %v", n, t)))
	} else {
		return n
	}
}

// construct a syntax error for unexpected token at current position.
func (p *Parser) Unexpected(t Token) *SyntaxError {
	return p.SyntaxError(fmt.Sprintf("unexpected %v", t))
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
