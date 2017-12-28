package e

import (
	"bytes"
	"fmt"
	"io"
)

type Expr interface {
	PrintTo(w io.Writer)
}

func String(e Expr) string {
	var buf bytes.Buffer
	e.PrintTo(&buf)
	return buf.String()
}

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

type Call struct {
	Func Expr
	Args []Expr
}

func (n *Call) PrintTo(w io.Writer) {
	n.Func.PrintTo(w)
	fmt.Fprint(w, "(")
	for i, a := range n.Args {
		if i != 0 {
			fmt.Fprint(w, ",")
		}
		a.PrintTo(w)
	}
	fmt.Fprint(w, ")")
}
