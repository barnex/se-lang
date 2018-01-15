package eva

import "github.com/barnex/se-lang/ast"

func compileLambda(n *ast.Lambda) Prog {
	panic("todo")
}

type Lambda struct {
}

func (n *Lambda) Eval() Value {
	panic("todo")
}

//func (n *PLambda) Apply(args []Value) Value {
//	for i, a := range n.Args {
//		a.Push(args[i])
//	}
//	v := n.Body.Eval()
//	for _, a := range n.Args {
//		a.Pop()
//	}
//	return v
//}
