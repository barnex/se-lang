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
	p := &LambdaProg{
		Body: compileExpr(n.Body),
	}
	for _, c := range n.Caps {
		p.Caps = append(p.Caps, compileVar(c.Src))
	}
	return p
}

type LambdaProg struct {
	Caps []Prog
	Body Prog
}

func (p *LambdaProg) Exec(m *Machine) {
	v := LambdaValue{Body: p.Body}
	for _, c := range p.Caps {
		c.Exec(m)
		v.Capv = append(v.Capv, m.RA())
	}
	m.SetRA(&v)
}

type LambdaValue struct {
	Capv []Value
	Body Prog
}

var _ Applier = (*LambdaValue)(nil)

func (p *LambdaValue) Apply(m *Machine) {
	m.Push(m.BP())
	m.SetBP(m.SP())
	for _, c := range p.Capv {
		m.Push(c)
	}
	p.Body.Exec(m)
	m.Grow(-len(p.Capv))
	m.SetBP(m.Pop().(int))
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
	if id.Var == nil {
		return compileGlobal(id)
	} else {
		return compileVar(id.Var)
	}
}

func compileGlobal(id *ast.Ident) Prog {
	p := prelude.Find(id.Name)
	if p == nil {
		panic(se.Errorf("compileIdent: undefined: %q: %#v", id.Name, id))
	}
	return p
}

func compileVar(v ast.Var) Prog {
	switch v := v.(type) {
	default:
		panic(unhandled(v))
	case nil:
		panic(unhandled(v))
	case *ast.Arg:
		return compileArg(v)
	case *ast.LocVar:
		return compileLocVar(v)
	}
}

func compileArg(a *ast.Arg) Prog {
	return fromBP{Offset: -2 - a.Index}
}

func compileLocVar(a *ast.LocVar) Prog {
	return fromBP{Offset: a.Index}
}

type fromBP struct {
	Offset int
}

func (p fromBP) Exec(m *Machine) {
	m.SetRA(m.FromBP(p.Offset))
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
