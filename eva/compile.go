package eva

import (
	"fmt"
	"strconv"

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
	case *ast.Block:
		return compileBlock(n)
	case *ast.Call:
		return compileCall(n)
	case *ast.Cond:
		return compileCond(n)
	case *ast.Ident:
		return compileIdent(n)
	case *ast.Lambda:
		return compileLambda(n)
	case *ast.Num:
		return compileNum(n)
	}
}

// -------- Block

func compileBlock(n *ast.Block) Prog {
	b := &Block{}
	for _, stmt := range n.Stmts {
		if a, ok := stmt.(*ast.Assign); ok {
			b.Init = append(b.Init, compileAssign(a))
		} else {
			if b.Expr != nil {
				panic(se.Errorf("block has more than 1 expression"))
			}
			b.Expr = compileExpr(stmt)
		}
	}
	if b.Expr == nil {
		panic(se.Errorf("block has no expression"))
	}
	return b
}

type Block struct {
	Init []Assign
	Expr Prog
}

func (b *Block) Exec(m *Machine) {
	for _, ini := range b.Init {
		ini.Exec(m)
	}
	b.Expr.Exec(m)
}

type Assign struct {
	LHS fromBP
	RHS Prog
}

func compileAssign(n *ast.Assign) Assign {
	return Assign{
		LHS: compileLocVar(n.LHS.Var.(*ast.LocVar)),
		RHS: compileExpr(n.RHS),
	}
}

func (a Assign) Exec(m *Machine) {
	a.RHS.Exec(m)
	a.LHS.SetToRA(m)
}

// -------- Cond

type Cond struct {
	Test, If, Else Prog
}

func compileCond(n *ast.Cond) *Cond {
	return &Cond{
		Test: compileExpr(n.Test),
		If:   compileExpr(n.If),
		Else: compileExpr(n.Else),
	}
}

func (p *Cond) Exec(m *Machine) {
	p.Test.Exec(m)
	if m.RA().Get().(bool) {
		p.If.Exec(m)
	} else {
		p.Else.Exec(m)
	}
}

// -------- Lambda

func compileLambda(n *ast.Lambda) Prog {
	p := &LambdaProg{
		Body:      compileExpr(n.Body),
		NumLocals: n.NumVar,
	}
	for _, c := range n.Caps {
		p.Caps = append(p.Caps, compileVar(c.Src))
	}
	return p
}

type LambdaProg struct {
	Caps      []Prog
	Body      Prog
	NumLocals int
}

func (p *LambdaProg) Exec(m *Machine) {
	v := &LambdaValue{Body: p.Body, NumLocals: p.NumLocals}
	for _, c := range p.Caps {
		c.Exec(m)
		if m.RA().v == nil {
			panic("capv==nil")
		}
		v.Capv = append(v.Capv, m.RA())
	}
	m.SetRA(box(v))
}

type LambdaValue struct {
	Capv      []Box
	Body      Prog
	NumLocals int
}

var _ Applier = (*LambdaValue)(nil)

func (p *LambdaValue) Apply(m *Machine) {
	m.Push(box(m.BP()))
	m.SetBP(m.SP())
	m.Grow(p.NumLocals)
	for i, c := range p.Capv {
		m.FromBP(i).Set(c.Get())
	}
	p.Body.Exec(m)
	m.Grow(-p.NumLocals)
	m.SetBP(m.Pop().Get().(int))
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
	p.F.Exec(m)                     // eval the function
	m.RA().Get().(Applier).Apply(m) // apply function to arguments
	m.Grow(-len(p.Args))            // free arguments stack space
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

func compileLocVar(a *ast.LocVar) fromBP {
	return fromBP{Offset: a.Index}
}

type fromBP struct {
	Offset int
}

func (p fromBP) Exec(m *Machine) {
	m.SetRA(m.FromBP(p.Offset))
}

func (p fromBP) SetToRA(m *Machine) {
	m.FromBP(p.Offset).Set(m.RA().Get())
}

// -------- Const

type Const struct {
	v Value
}

func (c Const) Exec(m *Machine) {
	m.SetRA(box(c.v))
}

func compileNum(n *ast.Num) Prog {
	if v, err := strconv.Atoi(n.Value); err == nil {
		return &Const{v} // int
	}
	v, err := strconv.ParseFloat(n.Value, 64)
	if err != nil {
		panic(se.Errorf("%v", err))
	}
	return Const{v}
}

// --------

func unhandled(x interface{}) string {
	return fmt.Sprintf("BUG: unhandled case: %T", x)
}
