package eva

import (
	"github.com/barnex/se-lang/ast"
)

func compileLambda(n *ast.Lambda) Prog {
	return &Lambda{
		NArgs: len(n.Args),
		Body:  compileExpr(n.Body),
	}
}

type Lambda struct {
	NArgs int
	Body  Prog
}

func (p *Lambda) Eval(s *Machine) {
	s.Push(p, "lambda: self")
}

func (p *Lambda) Apply(m *Machine) {
	p.Body.Eval(m)
	//	m.Push(666.6, "l")
}
