package main

import "fmt"

type Expr interface {
	String() string
}

type Num struct {
	Value float64
}

func (n *Num) String() string {
	return fmt.Sprint(n.Value)
}

type Ident struct {
	Name string
}

func (n *Ident) String() string {
	return n.Name
}

type Call struct {
	Func Expr
	Args []Expr
}

func (n *Call) String() string {
	str := n.Func.String() + "("
	for i, a := range n.Args {
		if i != 0 {
			str += ","
		}
		str += a.String()
	}
	str += ")"
	return str
}
