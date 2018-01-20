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
		Body: compileExpr(n.Body),
	}
	//p.Caps = make([]Capture, len(n.Caps))
	//for i := range p.Caps {
	//	p.Caps[i].Src = compileVar(n.Caps[i].Src)
	//	p.Caps[i].Dst = compileVar(n.Caps[i].Dst)
	//}
	return p
}

type Lambda struct {
	Caps []Capture
	Capv []Value
	Body Prog
}

type Capture struct {
	Src, Dst Var
}

var _ Applier = (*Lambda)(nil)

func (p_ *Lambda) Exec(m *Machine) {
	p := *p_
	// TODO: push captures
	//p.Capv = make([]Value, len(p.Caps))
	//for i := range p.Caps {
	//	p.Caps[i].Exec(m)
	//	//g
	//}
	m.SetRA(&p)
}

func (p *Lambda) Apply(m *Machine) {
	//m.Grow(len(p.Caps))
	//for i := range p.Caps {
	//	//m.s[m.BP+p.NArgs+i] = p.Caps[i]
	//	m.Push(p.Caps[i], "arg") // todo reverse order
	//}

	m.Push(m.BP())
	m.SetBP(m.SP())
	for _, c := range p.Capv {
		m.Push(c)
	}
	p.Body.Exec(m)
	m.Grow(-len(p.Capv))
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
	switch n := id.Var.(type) {
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
