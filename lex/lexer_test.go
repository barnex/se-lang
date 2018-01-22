package lex

import (
	"reflect"
	"strings"
	"testing"

	se "github.com/barnex/se-lang"
)

func TestLex(t *testing.T) {
	cases := []struct {
		src  string
		want []Token
	}{
		{``, []Token{}},
		{`//comment`, []Token{}},
		{"+", []Token{{TAdd, "+"}}},
		{"=", []Token{{TAssign, "="}}},
		{"/", []Token{{TDiv, "/"}}},
		{"==", []Token{{TEq, "=="}}},
		{"!=", []Token{{TNEq, "!="}}},
		{"123.4", []Token{{TNum, "123.4"}}},
		{">=", []Token{{TGe, ">="}}},
		{">", []Token{{TGt, ">"}}},
		{"ident", []Token{{TIdent, "ident"}}},
		{"1234", []Token{{TNum, "1234"}}},
		{"{", []Token{{TLBrace, "{"}}},
		{"(", []Token{{TLParen, "("}}},
		{"->", []Token{{TLambda, "->"}}},
		{"<=", []Token{{TLe, "<="}}},
		{"<", []Token{{TLt, "<"}}},
		{"-", []Token{{TMinus, "-"}}},
		{"*", []Token{{TMul, "*"}}},
		{"}", []Token{{TRBrace, "}"}}},
		{")", []Token{{TRParen, ")"}}},
		{`1`, []Token{{TNum, "1"}}},
		{`23`, []Token{{TNum, "23"}}},
		{` 45 	678 `, []Token{{TNum, "45"}, {TNum, "678"}}},
		{`x foo bar2`, []Token{{TIdent, "x"}, {TIdent, "foo"}, {TIdent, "bar2"}}},
		{` x foo bar0 `, []Token{{TIdent, "x"}, {TIdent, "foo"}, {TIdent, "bar0"}}},
		{`((foo )`, []Token{{TLParen, "("}, {TLParen, "("}, {TIdent, "foo"}, {TRParen, ")"}}},
		{` " a 1 () "`, []Token{{TString, `" a 1 () "`}}},
		{`""`, []Token{{TString, `""`}}},
		{`a+b*c`, []Token{{TIdent, "a"}, {TAdd, "+"}, {TIdent, "b"}, {TMul, "*"}, {TIdent, "c"}}},
		{`a==b`, []Token{{TIdent, "a"}, {TEq, "=="}, {TIdent, "b"}}},
		{`%`, []Token{{TMod, "%"}}},
		{`a&&b||!c`, []Token{{TIdent, "a"}, {TAnd, "&&"}, {TIdent, "b"}, {TOr, "||"}, {TNot, "!"}, {TIdent, "c"}}},
		{`x=1;x`, []Token{{TIdent, "x"}, {TAssign, "="}, {TNum, "1"}, {TSemicol, ";"}, {TIdent, "x"}}},
		{`'x`, []Token{{TQuote, "'"}, {TIdent, "x"}}},
		{`x?1:2`, []Token{{TIdent, "x"}, {TQuestion, "?"}, {TNum, "1"}, {TColon, ":"}, {TNum, "2"}}},
	}

	for _, c := range cases {
		have, err := lexAll(c.src)
		if err != nil {
			t.Errorf("%v: error: %v", c.src, err)
			continue
		}
		want := append(c.want, Token{TEOF, ""})
		if !reflect.DeepEqual(have, want) {
			t.Errorf("%v: have %v, want %v", c.src, have, want)
		}
	}
}

func TestError(t *testing.T) {
	cases := []string{
		`$`,
		`"`,
		`"""`,
	}

	for _, src := range cases {
		_, err := lexAll(src)
		if err == nil {
			t.Errorf("%v: expected error", src)
			continue
		}
	}
}

// lexAll splits a string in tokens.
func lexAll(input string) (t []Token, e error) {
	// catch syntax errors
	defer func() {
		switch err := recover().(type) {
		default:
			panic(err) // resume
		case nil:
			// no error
		case se.Error:
			t = nil
			e = err
		}
	}()
	l := NewLexer(strings.NewReader(input))

	// read all tokens
	var out []Token
	tok := l.Next()
	for {
		out = append(out, tok)
		if tok.TType == TEOF {
			break
		}
		tok = l.Next()
	}

	return out, nil
}
