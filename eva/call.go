package eva

import "github.com/barnex/se-lang/ast"

func compileCall(n *ast.Call) Prog {
	args := make([]Prog, len(n.Args))
	for i, a := range n.Args {
		args[i] = compileExpr(a)
	}
	f := compileExpr(n.F) // todo message
	return &PCall{f, args}
}

type PCall struct {
	F    Prog
	Args []Prog
}

func (n *PCall) Eval() Value {
	args := make([]Value, len(n.Args))
	for i, a := range n.Args {
		args[i] = a.Eval()
	}
	return n.F.Eval().(Applier).Apply(args)
}

type Applier interface {
	Apply([]Value) Value
}
