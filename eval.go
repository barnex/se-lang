package e

import (
	"fmt"
	"io"
)

func EvalSafe(n Node) (Node, error) {
	return withCatch(func() Node {
		Resolve(&prelude, n)
		return eval(&Machine{}, n)
	})
}

func eval(m *Machine, n Node) Node {
	switch n := n.(type) {
	default:
		panic(fmt.Sprintf("bug: %T %v", n, n))
	case Func:
		return n
	case *Num:
		return n
	case *Ident:
		return n.Value // presumably Resolved before
	case Local:
		return m.Get(-(n.N + 1))
	case List:
		return evalList(m, n)
	}
}

func evalList(m *Machine, l List) Node {
	return eval(m, l.Car()).(Applier).Apply(m, l.Cdr())
}

//func add(e *Env, l List) Node {
//	sum := 0.0
//	for _, n := range l {
//		sum += eval(e, n).(*Num).Value
//	}
//	return &Num{sum}
//}
//
//func mul(e *Env, l List) Node {
//	prod := 1.0
//	for _, n := range l {
//		prod *= eval(e, n).(*Num).Value
//	}
//	return &Num{prod}
//}
//
//func lambda(e *Env, l List) Node {
//	panic("todo")
//}
// _ Node = (Func)(nil)

type Applier interface {
	Apply(*Machine, List) Node
}

type Func func(*Machine, List) Node

func (f Func) PrintTo(w io.Writer) {
	fmt.Fprint(w, "func", f)
}

func (f Func) Apply(m *Machine, l List) Node {
	return f(m, l)
}
