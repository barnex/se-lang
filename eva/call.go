package eva

import (
	"fmt"

	"github.com/barnex/se-lang/ast"
)

func compileCall(n *ast.Call) Prog {
	f := compileExpr(n.F)
	args := make([]Prog, len(n.Args))
	for i, a := range n.Args {
		args[i] = compileExpr(a)
	}
	return &Call{f, args}
}

type Call struct {
	F    Prog
	Args []Prog
}

func (n *Call) Eval(m *Machine) {
	fmt.Println("call:", n.F)

	n.F.Eval(m)
	f := m.Pop("call:applier").(Applier)

	for i := len(n.Args) - 1; i >= 0; i-- {
		n.Args[i].Eval(m)
	}

	m.Push(m.EBP, "ebp")
	m.EBP = m.ESP()
	fmt.Println("ebp=", m.EBP)

	f.Apply(m)

	ret := m.Pop("call:return")

	m.EBP = m.Pop("ebp").(int)

	for range n.Args {
		m.Pop("call:shrink")
	}

	m.Push(ret, "call:return")
}

type Applier interface {
	Apply(s *Machine)
}
