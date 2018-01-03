package e

import (
	"fmt"
)

//func Eval(src string) (Node, error) {
//	fmt.Println("eval", src)
//	return withCatch(NewParser(strings.NewReader(src)).PExpr)
//}

func EvalNode(n Node) (Node, error) {
	return withCatch(func() Node {
		return eval(&scope, n)
	})
}

func eval(e *Env, n Node) Node {
	switch n := n.(type) {
	default:
		panic(fmt.Sprintf("bug: %T %v", n, n))
	case Func:
		return n
	case *Num:
		return n
	case *Ident:
		return e.Lookup(n.Name)
	case List:
		return evalList(e, n)
	}
}

var scope = Env{symbols: make(map[string]Node)}

type Env struct {
	parent  *Env
	symbols map[string]Node
}

func (e *Env) Lookup(name string) Node {
	if n, ok := e.symbols[name]; ok {
		return n
	}
	if e.parent == nil {
		panic(SyntaxErrorf("undefined: %v", name))
	}
	return e.parent.Lookup(name)
}

func (e *Env) Def(name string, value Node) {
	if _, ok := e.symbols[name]; ok {
		panic(SyntaxErrorf("already defined: %v", name))
	}
	e.symbols[name] = value
}

func init() {
	scope.Def("lisp", Func(evalList))
	scope.Def("add", Func(add))
	scope.Def("mul", Func(mul))
	scope.Def("lambda", Func(lambda))
}

func evalList(e *Env, l List) Node {
	return eval(e, l.Car()).(Applier).Apply(l.Cdr())
}

func add(l List) Node {
	sum := 0.0
	for _, n := range l {
		sum += eval(n).(*Num).Value
	}
	return &Num{sum}
}

func mul(l List) Node {
	prod := 1.0
	for _, n := range l {
		prod *= eval(n).(*Num).Value
	}
	return &Num{prod}
}

func lambda(l List) Node {

}
