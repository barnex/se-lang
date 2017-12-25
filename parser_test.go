package main

import (
	"reflect"
	"testing"
)

func TestParseExpr(t *testing.T) {

	cases := []struct {
		in   string
		want Expr
	}{
		{`1`, num(1)},
		{` 1 `, num(1)},
		{`(1)`, num(1)},
		{` ( 1 ) `, num(1)},
		{`f`, ident("f")},
		{`f()`, call(ident("f"))},
		{`f(x)`, call(ident("f"), ident("x"))},
		{`f((x))`, call(ident("f"), ident("x"))},
		{`(f)(x)`, call(ident("f"), ident("x"))},
		{`f(x)(y)`, call(call(ident("f"), ident("x")), ident("y"))},
		{`f(x,y)`, call(ident("f"), ident("x"), ident("y"))},
		{`f(x,y,)`, call(ident("f"), ident("x"), ident("y"))},
	}

	for _, c := range cases {
		have, err := Parse(c.in)
		if err != nil {
			t.Errorf("%v: error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

func TestParseToString(t *testing.T) {

	cases := []struct {
		in   string
		want string
	}{
		{` 1 `, `1`},
		{` (1) `, `1`},
		{`f`, `f`},
		{`f()`, `f()`},
		{`f(x)`, `f(x)`},
		{`f((x))`, `f(x)`},
		{`(f)(x)`, `f(x)`},
		{`f(x)(y)`, `f(x)(y)`},
		{`(f)(x,y,)`, `f(x,y)`},
	}

	for _, c := range cases {
		have, err := Parse(c.in)
		if err != nil {
			t.Errorf("%v: error: %v", c.in, err)
			continue
		}
		if have := have.String(); have != c.want {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

func TestParseError(t *testing.T) {
	cases := []string{
		`(1`,
		`1)`,
		` ( 1 `,
		` 1 ) `,
		`f(x`,
		`f(x))`,
		`f(x y)`,
		`f(,)`,
		`f g`,
		`f(g) x`,
		`1 2`,
	}

	for _, c := range cases {
		_, err := Parse(c)
		if err == nil {
			t.Errorf("%v: expected error", c)
			continue
		}
	}
}

func num(v float64) Expr             { return &Num{v} }
func ident(n string) Expr            { return &Ident{n} }
func call(f Expr, args ...Expr) Expr { return &Call{f, args} }
