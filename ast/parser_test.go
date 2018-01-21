package ast

import (
	"reflect"
	"strings"
	"testing"
)

// Parse expressions and compare to the expected AST.
func TestParseExpr(t *testing.T) {

	var (
		add = ident("add")
		f   = ident("f")
		mul = ident("mul")
		neg = ident("neg")
		one = num(1)
		sub = ident("sub")
		x   = ident("x")
		y   = ident("y")
		z   = ident("z")
	)

	cases := []struct {
		in   string
		want Node
	}{
		// operand
		//  | - operand
		{`--1`, call(neg, call(neg, one))},
		{`-1`, call(neg, one)},
		{`-f`, call(neg, f)},
		{`-(f)`, call(neg, f)},
		{`!x`, call(ident("not"), x)},

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

		// binary
		{`1*2>3`, call(ident("gt"), call(mul, num(1), num(2)), num(3))},
		{`1*2<3`, call(ident("lt"), call(mul, num(1), num(2)), num(3))},
		{`1*2>=3`, call(ident("ge"), call(mul, num(1), num(2)), num(3))},
		{`1*2<=3`, call(ident("le"), call(mul, num(1), num(2)), num(3))},
		{`1*2==3`, call(ident("eq"), call(mul, num(1), num(2)), num(3))},
		{`1*2!=3`, call(ident("neq"), call(mul, num(1), num(2)), num(3))},
		{`2-1`, call(sub, num(2), num(1))},
		//{`3%4`, call(ident("mod"), num(3), num(4))},

		// random
		{`(f)(x)`, call(f, x)},
		{`f(x)(y)`, call(call(f, x), y)},
		{`1+2+3`, call(add, call(add, num(1), num(2)), num(3))},
		{`1+2*3`, call(add, num(1), call(mul, num(2), num(3)))},
		{`1*2+3`, call(add, call(mul, num(1), num(2)), num(3))},

		// lambda
		{`x->y`, lambda(args(x), y)},
		{`(x)->(x)`, lambda(args(x), x)},
		{`x->-y`, lambda(args(x), call(neg, y))},
		{`(x,y)->f(y,x)`, lambda(args(x, y), call(f, y, x))},
		{`(x,y)->f(y,x)()`, lambda(args(x, y), call(call(f, y, x)))},
		{`((x,y)->y+x)()`, call(lambda(args(x, y), call(add, y, x)))},
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

func num(v float64) Node                   { return &Num{v} }
func ident(n string) *Ident                { return &Ident{Name: n} }
func call(f Node, args ...Node) Node       { return &Call{f, normalize(args)} }
func lambda(args []*Ident, body Node) Node { return &Lambda{Args: args, Body: body} }
func args(n ...*Ident) []*Ident            { return n }

func normalize(x []Node) []Node {
	if x == nil {
		return []Node{}
		// reflect.DeepEqual considers nil different from empty list
	}
	return x
}
