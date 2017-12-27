package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {

	cases := []struct {
		in   string
		want []Token
	}{
		{``, []Token{{TEOF, ""}}},
		{`1`, []Token{{TNum, "1"}, {TEOF, ""}}},
		{`23`, []Token{{TNum, "23"}, {TEOF, ""}}},
		{` 45 	678 `, []Token{{TNum, "45"}, {TNum, "678"}, {TEOF, ""}}},
		{`x foo bar2`, []Token{{TIdent, "x"}, {TIdent, "foo"}, {TIdent, "bar2"}, {TEOF, ""}}},
		{` x foo bar0 `, []Token{{TIdent, "x"}, {TIdent, "foo"}, {TIdent, "bar0"}, {TEOF, ""}}},
		{`2x`, []Token{{TNum, "2x"}, {TEOF, ""}}}, // let atoi catch the syntax error
		{`((foo )`, []Token{{TLParen, "("}, {TLParen, "("}, {TIdent, "foo"}, {TRParen, ")"}, {TEOF, ""}}},
		{` " a 1 () "`, []Token{{TString, `" a 1 () "`}, {TEOF, ""}}},
		{`""`, []Token{{TString, `""`}, {TEOF, ""}}},
		{`a+b*c`, []Token{{TIdent, "a"}, {TAdd, "+"}, {TIdent, "b"}, {TMul, "*"}, {TIdent, "c"}, {TEOF, ""}}},
	}

	for _, c := range cases {
		have, err := Lex(c.in)
		if err != nil {
			t.Errorf("%v: error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

func TestError(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`$`, "illegal character"},
		{`"`, "unterminated string"},
		{`"""`, "unterminated string"},
	}

	for _, c := range cases {
		_, err := Lex(c.in)
		if err == nil {
			t.Errorf("%v: expected error", c.in)
			continue
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: have: %v, want: %v", c.in, err.Error(), c.want)
		}
	}
}
