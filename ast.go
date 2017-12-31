package e

import (
	"bytes"
	"fmt"
	"io"
)

type Node interface {
	PrintTo(w io.Writer)
}

var (
	_ Node = (*Num)(nil)
	_ Node = (*Ident)(nil)
	_ Node = (*List)(nil)
)

type Num struct {
	Value float64
}

func (n *Num) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Value)
}

type Ident struct {
	Name string
}

func (n *Ident) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Name)
}

type List []Node

func (l List) Car() Node {
	if len(l) == 0 {
		return nil
	}
	return l[0]
}

func (l List) Cdr() List {
	return l[1:]
}

func MakeList(car Node, cdr ...Node) List {
	a := make(List, 1+len(cdr))
	a[0] = car
	copy(a[1:], cdr)
	return a
}

func (n List) PrintTo(w io.Writer) {
	fmt.Fprint(w, "(")
	for i, a := range n {
		if i != 0 {
			fmt.Fprint(w, " ")
		}
		a.PrintTo(w)
	}
	fmt.Fprint(w, ")")
}

func ToString(e Node) string {
	var buf bytes.Buffer
	e.PrintTo(&buf)
	return buf.String()
}

// -------- Value
//type Value = reflect.Value
//
//
//var scope = map[string]reflect.Value{
//	"sqrt": reflect.ValueOf(math.Sqrt),
//	"+":    reflect.ValueOf(add),
//	"-":    reflect.ValueOf(sub),
//	"*":    reflect.ValueOf(mul),
//	"/":    reflect.ValueOf(div),
//}
//
//func add(x, y float64) float64 { return x + y }
//func sub(x, y float64) float64 { return x - y }
//func mul(x, y float64) float64 { return x * y }
//func div(x, y float64) float64 { return x / y }
//
//func (n *Ident) Eval() Value {
//	if v, ok := scope[n.Name]; ok {
//		return v
//	}
//	panic("undefined: " + n.Name)
//}
//
////func (n *Call) Eval() Value {
////	f := n.Car.Eval()
////	args := make([]Value, len(n.Cdr))
////	for i, a := range n.Cdr {
////		args[i] = a.Eval()
////	}
////	ret := f.Call(args)
////	if n := len(ret); n != 1 {
////		panic(fmt.Sprint(n, "return values"))
////	}
////	return ret[0]
////}
//
//func (n *List) Eval() Value {
//	panic("TODO")
//}
