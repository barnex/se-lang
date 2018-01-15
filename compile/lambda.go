package se

func compileLambda(n *Lambda) Prog {
	panic("todo")
}

type PLambda struct {
}

func (n *PLambda) Eval() Value {
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
