package se

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// A Node is an element of an AST (Abstract Syntax Tree).
type Node interface {
	PrintTo(w io.Writer)
}

// Num is a number Node, e.g.: '1'
type Num struct {
	Value float64
}

func (n *Num) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Value)
}

// Ident is an identifier Node, e.g.: 'sqrt'
type Ident struct {
	Name                string
	enclosing, defining *Scope
}

func (n *Ident) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Name)
	//if n.Stack.ID != 0 {
	//	fmt.Fprint(w, ":", n.ID)
	//}
}

// Call is a function call Node, e.g.: 'sqrt(2)'
type Call struct {
	F    Node
	Args []Node
}

func (n *Call) PrintTo(w io.Writer) {
	n.F.PrintTo(w)
	printList(w, n.Args)
}

// Lambda is a lambda expression node, e.g.: 'x->x*x'
type Lambda struct {
	Args  []*Ident
	Body  Node
	scope *Scope
}

func (n *Lambda) PrintTo(w io.Writer) {
	fmt.Fprint(w, "(")
	printList(w, n.Args)
	fmt.Fprint(w, TLambda)
	n.Body.PrintTo(w)
	fmt.Fprint(w, ")")
}

// printList prints a slice whose elements implement Node, e.g.:
//  []Node, []*Num, []*Ident, []*Call, []*Lambda
func printList(w io.Writer, list interface{}) {
	l := reflect.ValueOf(list)
	fmt.Fprint(w, "(")
	for i := 0; i < l.Len(); i++ {
		if i != 0 {
			fmt.Fprint(w, ", ")
		}
		l.Index(i).Interface().(Node).PrintTo(w)
	}
	fmt.Fprint(w, ")")
}

// ToString returns a string representation based on PrintTo
func ToString(e Node) string {
	var buf bytes.Buffer
	e.PrintTo(&buf)
	return buf.String()
}
