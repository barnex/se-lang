package e

import (
	"reflect"
	"strings"
	"testing"
)

// Parse expressions and compare to the expected AST.
func TestParseExpr(t *testing.T) {

	var (
		add    = ident("add")
		f      = ident("f")
		lambda = ident("lambda")
		mul    = ident("mul")
		neg    = ident("neg")
		one    = num(1)
		x      = ident("x")
		y      = ident("y")
		z      = ident("z")
	)

	cases := []struct {
		in   string
		want Node
	}{
		// operand
		//  | - operand
		{`--1`, list(neg, list(neg, one))},
		{`-1`, list(neg, one)},
		{`-f`, list(neg, f)},
		{`-(f)`, list(neg, f)},
		//  | num
		{`1`, one},
		//  | ident
		{`f`, f},
		//  | ( expr )
		{`(-1)`, list(neg, num(1))},
		{`(1)`, num(1)},
		{`(f)`, f},
		{`((f))`, f},
		//  | operand *(list)
		{`f()`, list(f)},
		{`f(x)`, list(f, x)},
		{`f(x,y,z)`, list(f, x, y, z)},
		{`(f)(x,y,z)()`, list(list(f, x, y, z))},

		// random
		{`(f)(x)`, list(f, x)},
		{`f(x)(y)`, list(list(f, x), y)},
		{`1+2+3`, list(add, list(add, num(1), num(2)), num(3))},
		{`1+2*3`, list(add, num(1), list(mul, num(2), num(3)))},
		{`1*2+3`, list(add, list(mul, num(1), num(2)), num(3))},

		// lambda
		{`x->y`, list(lambda, list(x), y)},
		{`x->-y`, list(lambda, list(x), list(neg, y))},
		{`(x,y)->f(y,x)`, list(lambda, list(x, y), list(f, y, x))},
		{`(x,y)->f(y,x)()`, list(lambda, list(x, y), list(list(f, y, x)))},
		{`((x,y)->y+x)()`, list(list(lambda, list(x, y), list(add, y, x)))},
	}

	for i, c := range cases {
		have, err := parse(c.in)
		if err != nil {
			t.Errorf("case %v: %v: error: %v", i, c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("case %v: %v: have %v, want %v", i, c.in, ToString(have), ToString(c.want))
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
		`+`,
		`-`,
		`*`,
		`,`,
		`(,)`,
		`1+`,
		`a-`,
		`(1+1)->2`,
		`x(y)->x+y`, // not (lambda (x y) (add x y))
		`()()`,      // not (())
		`()`,        // not ()
		`(1,2)`,     // not (1 2)
		`1,2`,       // not (1 2)
		`1,x`,
		`x,y->y,x`,
		`(x,y)->(y,x)`,
	}

	for _, c := range cases {
		e, err := parse(c)
		if err == nil {
			t.Errorf("%v: expected error, have: %v", c, ToString(e))
			continue
		}
	}
}

// Parse expressions and turn the AST's to strings.
func TestParseToString(t *testing.T) {

	cases := []struct {
		in   string
		want string
	}{
	//{` 1 `, `1`},
	//{` (1) `, `1`},
	//{`f`, `f`},
	//{`f()`, `(f)`},
	//{`f(x)`, `(f x)`},
	//{`f((x))`, `(f x)`},
	//{`(f)(x)`, `(f x)`},
	//{`f(x)(y)`, `((f x) y)`},
	////{`(f)(x,y)`, `(f (, x y))`},
	//{`x+y`, `(+ x y)`},
	//{`x*y`, `(* x y)`},
	//{`a+b*c`, `(+ a (* b c))`},
	//{`a*b+c`, `(+ (* a b) c)`},
	//{`a*(b+c)`, `(* a (+ b c))`},
	////{`x->x*x`, `(-> x (* x x))`},
	////{`sum=(x,y)->(x+y)`, `(= sum (-> (, x y) (+ x y)))`},
	////{`f=()->(3,4)`, `(= f (-> (,) (, 3 4)))`},
	//{`()`, `(list)`},
	}

	for _, c := range cases {
		have, err := parse(c.in)
		if err != nil {
			t.Errorf("%v: error: %v", c.in, err)
			continue
		}
		if have := ToString(have); have != c.want {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

func parse(src string) (Node, error) {
	return Parse(strings.NewReader(src))
}
