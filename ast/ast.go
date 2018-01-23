package ast

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/barnex/se-lang/lex"
)

// A Node is an element of an AST (Abstract Syntax Tree).
type Node interface {
	PrintTo(w io.Writer)
}

// Assign is a declaration, e.g.: a=1
type Assign struct {
	LHS *Ident
	RHS Node
}

func (n *Assign) PrintTo(w io.Writer) {
	n.LHS.PrintTo(w)
	fmt.Fprint(w, lex.TAssign)
	n.RHS.PrintTo(w)
}

// Block is a list of statements, e.g.: {a=1; b}
type Block struct {
	Stmts []Node
}

func (n *Block) PrintTo(w io.Writer) {
	fmt.Fprint(w, lex.TLBrace)
	for _, s := range n.Stmts {
		s.PrintTo(w)
		fmt.Fprint(w, lex.TSemicol)
	}
	fmt.Fprint(w, lex.TRBrace)
}

type Cond struct {
	Test, If, Else Node
}

func (n *Cond) PrintTo(w io.Writer) {
	n.Test.PrintTo(w)
	fmt.Fprint(w, lex.TQuestion)
	n.If.PrintTo(w)
	fmt.Fprint(w, lex.TColon)
	n.Else.PrintTo(w)
}

// Num is a number Node, e.g.: '1'
type Num struct {
	Value string
}

func (n *Num) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Value)
}

// Ident is an identifier Node, e.g.: 'sqrt'
type Ident struct {
	Name string
	Var  // filled in later by resolve
}

func (n *Ident) PrintTo(w io.Writer) {
	fmt.Fprint(w, n.Name)
	if n.Var != nil {
		fmt.Fprint(w, ":", n.Var)
	} else {
		fmt.Fprint(w, "??")
	}
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
	Args   []*Ident
	Caps   []Capture // filled in by resolve
	NumVar int
	Body   Node
}

type Capture struct {
	Name string
	Src  Var // variable being captured from the parent frame
	Dst  Var // variable being captured to
}

func (c *Capture) String() string {
	return fmt.Sprint(c.Dst, "=", c.Src)
}

func (n *Lambda) PrintTo(w io.Writer) {
	fmt.Fprint(w, "(")
	printList(w, n.Args)

	fmt.Fprint(w, lex.TLambda)
	if len(n.Caps) > 0 {
		fmt.Fprint(w, "[")
		for _, c := range n.Caps {
			fmt.Fprint(w, c, ",")
		}
		fmt.Fprint(w, "]")
	}

	n.Body.PrintTo(w)
	fmt.Fprint(w, ")")
}

func (n *Lambda) NewVariable() Var {
	v := &LocVar{Index: n.NumVar}
	n.NumVar++
	return v
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
