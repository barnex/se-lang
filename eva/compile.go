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
	f := compileExpr(n.F)
	args := make([]Prog, len(n.Args))
	for i, a := range n.Args {
		args[i] = compileExpr(a)
	}
	return &Call{f, args}
}

func (n *Call) Exec(m *Machine) {
	//fmt.Printf("eval %#v\n", n)
	//n.F.Eval(m)
	//f := m.RA.(Applier)

	////m.Grow(f.NFrame())

	//for i := len(n.Args) - 1; i >= 0; i-- {
	//	n.Args[i].Eval(m)
	//	//fmt.Println("stack bp+", i, "=", m.RA)
	//	//m.s[m.BP+i] = m.RA
	//	m.Push(m.RA, fmt.Sprint("arg", i))
	//}

	//m.Push(m.BP, "call-preamble")
	//m.BP = m.SP()
	//fmt.Println("bp=", m.BP)

	//f.Apply(m)

	//m.BP = m.Pop("call-restore-bp").(int)
	//m.Grow(-f.NFrame())
}

//type Applier interface {
//	Apply(s *Machine)
//	NFrame() int
//}

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
