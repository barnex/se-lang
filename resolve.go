package e

//func Resolve(s *Scope, n Node) {
//	switch n := n.(type) {
//	default:
//		return
//	case *Ident:
//		n.Value = s.Resolve(n.Name)
//	case List:
//		for _, n := range n {
//			Resolve(s, n)
//		}
//	}
//}
//
//type Scope struct {
//	parent  *Scope
//	symbols map[string]Node
//}
//
//func (e *Scope) Resolve(name string) Node {
//	if n, ok := e.symbols[name]; ok {
//		return n
//	}
//	if e.parent == nil {
//		panic(SyntaxErrorf("undefined: %v", name))
//	}
//	return e.parent.Resolve(name)
//}
//
//func (e *Scope) Def(name string, value Node) {
//	if _, ok := e.symbols[name]; ok {
//		panic(SyntaxErrorf("already defined: %v", name))
//	}
//	e.symbols[name] = value
//}
//
//var prelude = Scope{symbols: make(map[string]Node)}
//
//func init() {
//	prelude.Def("pi", num(math.Pi))
//	//prelude.Def("add", Func(add))
//	//prelude.Def("mul", Func(mul))
//	//prelude.Def("lambda", Func(lambda))
//	//prelude.Def("lisp", Func(evalList))
//}
