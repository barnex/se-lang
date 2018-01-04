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
	_ Node = (List)(nil)
	_ Node = Local{}
)

type Num struct {
	Value float64
}

func (n *Num) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Value)
}

type Ident struct {
	Name  string
	Value Node // if name is resolved
}

func (n *Ident) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Name)
}

type Local struct {
	N int
}

func (n Local) PrintTo(w io.Writer) {
	fmt.Fprint(w, "local", n.N)
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
