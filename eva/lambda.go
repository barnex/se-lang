package eva

import (
	"github.com/barnex/se-lang/ast"
)

func compileLambda(n *ast.Lambda) Prog {

	// TODO: captured arguments should come before regular args, not after

	ncaps := make([]int, len(n.Cap))
	for i := range ncaps {
		ncaps[i] = n.Cap[i].ParVar.Index
	}
	return &Lambda{
		NArgs: len(n.Args),
		NCaps: ncaps,
		Body:  compileExpr(n.Body),
	}
}

type Lambda struct {
	NArgs int
	NCaps []int
	Caps  []Value
	Body  Prog
}

func (p_ *Lambda) Eval(s *Machine) {
	cpy := *p_
	cpy.Caps = make([]Value, len(cpy.NCaps))
	for i := range cpy.NCaps {
		cpy.Caps[i] = s.FromBP(cpy.NCaps[i], "capture")
	}
	//fmt.Printf("lambda: eval: self=%#v\n", cpy)
	s.RA = &cpy
}

func (p *Lambda) Apply(m *Machine) {
	//m.Grow(len(p.Caps))
	for i := range p.Caps {
		//m.s[m.BP+p.NArgs+i] = p.Caps[i]
		m.Push(p.Caps[i], "arg") // todo reverse order
	}

	p.Body.Eval(m)

	//if len(p.NCaps) > 0 {
	//	// free captures
	//	ret := m.Pop("lambda:sub-ret")
	//	for i := len(p.NCaps) - 1; i >= 0; i-- {
	//		m.Pop("lambda:free-cap")
	//	}
	//	m.Push(ret, "lambda:return")
	//}
}

func (p *Lambda) NFrame() int {
	return p.NArgs + len(p.NCaps)
}

var _ Applier = (*Lambda)(nil)
