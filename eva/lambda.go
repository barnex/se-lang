package eva

import (
	"github.com/barnex/se-lang/ast"
)

func compileLambda(n *ast.Lambda) Prog {
	p := &Lambda{
		Args: n.Args,
		Caps: n.Caps,
		Body: compileExpr(n.Body),
	}
	for i := range p.Caps{
		p.Caps[i].Dst = &ast.
	}
}

type Lambda struct {
	Args []*ast.Ident
	Caps []*ast.CaptVar
	Capv []Value
	Body Prog
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
