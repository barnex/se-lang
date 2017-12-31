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
		want Expr
	}{
		// operand
		//  | - operand
		{`--1`, call(neg, call(neg, one))},
		{`-1`, call(neg, one)},
		{`-f`, call(neg, f)},
		{`-(f)`, call(neg, f)},
		//  | num
		{`1`, one},
		//  | ident
		{`f`, f},
		//  | ( expr )
		{`(-1)`, call(neg, num(1))},
		{`(1)`, num(1)},
		{`(f)`, f},
		{`((f))`, f},
		//  | operand *(list)
		{`f()`, call(f)},
		{`f(x)`, call(f, x)},
		{`f(x,y,z)`, call(f, x, y, z)},
		{`(f)(x,y,z)()`, call(call(f, x, y, z))},

		// expr
		{`1, x`, list(num(1), x)},
		{`(1, x)`, list(num(1), x)},
		{`()`, list()},

		// random
		{`(f)(x)`, call(f, x)},
		{`f(x)(y)`, call(call(f, x), y)},
		{`1+2+3`, call(add, call(add, num(1), num(2)), num(3))},
		{`1+2*3`, call(add, num(1), call(mul, num(2), num(3)))},
		{`1*2+3`, call(add, call(mul, num(1), num(2)), num(3))},
		{`(x,y)`, list(x, y)},

		// lambda
		{`x->y`, call(lambda, list(x), y)},
		{`x->-y`, call(lambda, list(x), call(neg, y))},
		{`x,y->y,x`, list(x, call(lambda, list(y), y), x)},
		{`(x,y)->(y,x)`, call(lambda, list(x, y), list(y, x))},
		{`(x,y)->f(y,x)`, call(lambda, list(x, y), call(f, y, x))},
		{`(x,y)->f(y,x)()`, call(lambda, list(x, y), call(call(f, y, x)))},
		{`((x,y)->y+x)()`, call(call(lambda, list(x, y), call(add, y, x)))},
	}

	for i, c := range cases {
		have, err := parse(c.in)
		if err != nil {
			t.Errorf("case %v: %v: error: %v", i, c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("case %v: %v: have %v, want %v", i, c.in, ExprString(have), ExprString(c.want))
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
	}

	for _, c := range cases {
		e, err := parse(c)
		if err == nil {
			t.Errorf("%v: expected error, have: %v", c, ExprString(e))
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
		if have := ExprString(have); have != c.want {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

func parse(src string) (Expr, error) {
	return Parse(strings.NewReader(src))
}

func num(v float64) Expr             { return &Num{v} }
func ident(n string) Expr            { return &Ident{n} }
func call(f Expr, args ...Expr) Expr { return &Call{f, normalize(args)} }
func list(x ...Expr) Expr            { return &List{normalize(x)} }

func normalize(x []Expr) []Expr {
	if x == nil {
		return []Expr{}
		// reflect.DeepEqual considers nil different from empty list
	}
	return x
}
