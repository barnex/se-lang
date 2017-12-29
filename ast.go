package e

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
)

// ---- Expr

type Expr interface {
	Eval() Value
	PrintTo(w io.Writer)
}

func ExprString(e Expr) string {
	var buf bytes.Buffer
	e.PrintTo(&buf)
	return buf.String()
}

// -------- Atomic
// ---- Num

type Num struct {
	Value float64
}

func (n *Num) Eval() Value {
	return reflect.ValueOf(n.Value)
}

func (n *Num) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Value)
}

// ---- Ident

type Ident struct {
	Name string
}

var scope = map[string]reflect.Value{
	"sqrt": reflect.ValueOf(math.Sqrt),
	"+":    reflect.ValueOf(add),
	"-":    reflect.ValueOf(sub),
	"*":    reflect.ValueOf(mul),
	"/":    reflect.ValueOf(div),
}

func add(x, y float64) float64 { return x + y }
func sub(x, y float64) float64 { return x - y }
func mul(x, y float64) float64 { return x * y }
func div(x, y float64) float64 { return x / y }

func (n *Ident) Eval() Value {
	if v, ok := scope[n.Name]; ok {
		return v
	}
	panic("undefined: " + n.Name)
}

func (n *Ident) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Name)
}

// -------- Composite

type Comp struct {
	Car Expr
	Cdr []Expr
}

func (n *Comp) Eval() Value {
	f := n.Car.Eval()
	args := make([]Value, len(n.Cdr))
	for i, a := range n.Cdr {
		args[i] = a.Eval()
	}
	ret := f.Call(args)
	if n := len(ret); n != 1 {
		panic(fmt.Sprint(n, "return values"))
	}
	return ret[0]
}

func (n *Comp) PrintTo(w io.Writer) {
	fmt.Fprint(w, "(")
	n.Car.PrintTo(w)
	for _, a := range n.Cdr {
		//if i != 0 {
		fmt.Fprint(w, " ")
		//}
		a.PrintTo(w)
	}
	fmt.Fprint(w, ")")
}

// -------- Value
type Value = reflect.Value
