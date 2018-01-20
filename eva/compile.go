package eva

import (
	"fmt"

	"github.com/barnex/se-lang/ast"
)

type Prog interface {
	Exec(s *Machine)
}

func compileExpr(n ast.Node) Prog {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *ast.Call:
		return compileCall(n)
	case *ast.Ident:
		return compileIdent(n)
	case *ast.Lambda:
		return compileLambda(n)
	case *ast.Num:
		return compileNum(n)
	}
}

// -------- Call

type Call struct {
	F    Prog
	Args []Prog
}

func compileCall(n *ast.Call) Prog {
	var c Call
	c.F = compileExpr(n.F)
	for _, a := range n.Args {
		c.Args = append(c.Args, compileExpr(a))
	}
	return &c
}

func (p *Call) Exec(m *Machine) {
	for i := len(p.Args) - 1; i >= 0; i-- {
		p.Args[i].Exec(m) // eval argument
		m.Push(m.RA)      // push argument
	}
	p.F.Exec(m)             // eval the function
	m.RA.(Applier).Apply(m) // apply function to arguments
	m.Grow(-len(p.Args))    // free arguments stack space
}

type Applier interface {
	Apply(s *Machine)
}

// -------- Const

type Const struct {
	v Value
}

func (c Const) Exec(m *Machine) {
	m.RA = c.v
}

func compileNum(n *ast.Num) Prog {
	return &Const{n.Value}
}

// --------

func unhandled(x interface{}) string {
	return fmt.Sprintf("BUG: unhandled case: %T", x)
}
