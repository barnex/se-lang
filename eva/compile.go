package eva

import (
	"fmt"

	se "github.com/barnex/se-lang"
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

// -------- Lambda

func compileLambda(n *ast.Lambda) Prog {
	p := &Lambda{
		//Args: n.Args,
		//Caps: n.Caps,
		Body: compileExpr(n.Body),
	}
	//for i := range p.Caps{
	//	p.Caps[i].Dst = &ast.
	//}
	return p
}

type Lambda struct {
	//Args []*Arg
	//Caps []*ast.CaptVar
	//Capv []Value
	Body Prog
}

var _ Applier = (*Lambda)(nil)

func (p *Lambda) Exec(m *Machine) {
	//cpy := *p_
	//cpy.Caps = make([]Value, len(cpy.NCaps))
	//for i := range cpy.NCaps {
	//	cpy.Caps[i] = s.FromBP(cpy.NCaps[i], "capture")
	//}
	////fmt.Printf("lambda: eval: self=%#v\n", cpy)
	//s.RA = &cpy

	m.SetRA(p)
}

func (p *Lambda) Apply(m *Machine) {
	//m.Grow(len(p.Caps))
	//for i := range p.Caps {
	//	//m.s[m.BP+p.NArgs+i] = p.Caps[i]
	//	m.Push(p.Caps[i], "arg") // todo reverse order
	//}

	m.Push(m.BP())
	m.SetBP(m.SP())
	p.Body.Exec(m)
	m.SetBP(m.Pop().(int))

	//if len(p.NCaps) > 0 {
	//	// free captures
	//	ret := m.Pop("lambda:sub-ret")
	//	for i := len(p.NCaps) - 1; i >= 0; i-- {
	//		m.Pop("lambda:free-cap")
	//	}
	//	m.Push(ret, "lambda:return")
	//}
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
		m.Push(m.RA())    // push argument
	}
	p.F.Exec(m)               // eval the function
	m.RA().(Applier).Apply(m) // apply function to arguments
	m.Grow(-len(p.Args))      // free arguments stack space
}

type Applier interface {
	Apply(s *Machine)
}

// -------- Ident

func compileIdent(id *ast.Ident) Prog {
	switch n := id.Object.(type) {
	default:
		panic(unhandled(n))
	case nil:
		panic(se.Errorf("compileIdent: undefined: %q: %#v", id.Name, id))
	case Prog:
		return n
	}
}

// -------- Const

type Const struct {
	v Value
}

func (c Const) Exec(m *Machine) {
	m.SetRA(c.v)
}

func compileNum(n *ast.Num) Prog {
	return &Const{n.Value}
}

// --------

func unhandled(x interface{}) string {
	return fmt.Sprintf("BUG: unhandled case: %T", x)
}
