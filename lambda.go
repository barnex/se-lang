package se

import (
	"fmt"
	"io"
)

type Lambda struct {
	Args  []*Ident
	Body  Node
	scope *Scope
}

func (n *Lambda) PrintTo(w io.Writer) {
	printIdents(w, n.Args)
	fmt.Fprint(w, TLambda)
	n.Body.PrintTo(w)
}
