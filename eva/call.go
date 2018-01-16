package eva

import "github.com/barnex/se-lang/ast"

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

func (n *Call) Eval(s *Stack) {
	for _, a := range n.Args {
		a.Eval(s)
	}
	n.F.Eval(s)
	s.Pop().(Applier).Apply(s)
}

type Applier interface {
	Apply(s *Stack)
}
