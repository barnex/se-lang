package e

import (
	"reflect"
	"strings"
	"testing"
)

// Parse expressions and compare to the expected AST.
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
		{`1+2+3`, call(ident("+"), call(ident("+"), num(1), num(2)), num(3))},
	}

	for i, c := range cases {
		have, err := parse(c.in)
		if err != nil {
			t.Errorf("case %v: %v: error: %v", i, c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

// Parse expressions and turn the AST's to strings.
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
		{`x+y`, `+(x,y)`},
		{`x*y`, `*(x,y)`},
		{`a+b*c`, `+(a,*(b,c))`},
		{`a*b+c`, `+(*(a,b),c)`},
		{`a*(b+c)`, `*(a,+(b,c))`},
	}

	for _, c := range cases {
		have, err := parse(c.in)
		if err != nil {
			t.Errorf("%v: error: %v", c.in, err)
			continue
		}
		if have := String(have); have != c.want {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

// Ensure parse errors on bad input.
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
		e, err := parse(c)
		if err == nil {
			t.Errorf("%v: expected error, have: %v", c, String(e))
			continue
		}
	}
}

func parse(src string) (Expr, error) {
	return Parse(strings.NewReader(src))
}

func num(v float64) Expr             { return &Num{v} }
func ident(n string) Expr            { return &Ident{n} }
func call(f Expr, args ...Expr) Expr { return &Call{f, args} }
