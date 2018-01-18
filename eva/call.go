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
	fmt.Printf("eval %#v\n", n)
	n.F.Eval(m)
	f := m.RA.(Applier)

	//m.Grow(f.NFrame())

	for i := len(n.Args) - 1; i >= 0; i-- {
		n.Args[i].Eval(m)
		//fmt.Println("stack bp+", i, "=", m.RA)
		//m.s[m.BP+i] = m.RA
		m.Push(m.RA, fmt.Sprint("arg", i))
	}

	m.Push(m.BP, "call-preamble")
	m.BP = m.SP()
	fmt.Println("bp=", m.BP)

	f.Apply(m)

	m.BP = m.Pop("call-restore-bp").(int)
	m.Grow(-f.NFrame())

}

type Applier interface {
	Apply(s *Machine)
	NFrame() int
}
