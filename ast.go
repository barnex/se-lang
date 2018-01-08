package se

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
	_ Node = (*Call)(nil)
	_ Node = (*Lambda)(nil)
)

type Num struct {
	Value float64
}

func (n *Num) Eval() Value {
	return n.Value
}

func (n *Num) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Value)
}

type Ident struct {
	Name string
	Prog
}

func (n *Ident) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Name)
}

type Call struct {
	F    Node
	Args []Node
}

func (n *Call) Eval() Value {
	panic("todo: Call.Eval")
}

func (n *Call) PrintTo(w io.Writer) {
	n.F.PrintTo(w)
	printList(w, n.Args)
}

func printList(w io.Writer, l []Node) {
	fmt.Fprint(w, "(")
	for i, a := range l {
		if i != 0 {
			fmt.Fprint(w, ", ")
		}
		a.PrintTo(w)
	}
	fmt.Fprint(w, ")")
}

func printIdents(w io.Writer, l []*Ident) {
	fmt.Fprint(w, "(")
	for i, a := range l {
		if i != 0 {
			fmt.Fprint(w, ", ")
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
