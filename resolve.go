package se

import "fmt"

func Resolve(n Node) {
	resolve(Frames{prelude}, n)
}

func resolve(s Frames, n Node) {
	switch n := n.(type) {
	case *Ident:
		resolveIdent(s, n)
	case *Lambda:
		resolveLambda(s, n)
	case *Num, *Call:
		for _, n := range children(n) {
			resolve(s, n)
		}
	default:
		panic(unhandled(n))
	}
}

func resolveLambda(s Frames, n *Lambda) {
	// first define the arguments
	for i, a := range n.Args {
		a.Var = Local{i}
	}
	// then resolve the body
	s.Push(n)
	resolve(s, n.Body)
	s.Pop()
}

func resolveIdent(s Frames, id *Ident) {
	v, defScope := s.Find(id.Name)
	if v == nil {
		panic(SyntaxErrorf("undefined: %v", id.Name))
	}
	usingScope := s.Last()
	if defScope == usingScope {
		// local variable
		id.Var = v
	} else {
		// captured variable (TODO: or global).
		fmt.Println("capture", id.Name)
		usingScope.(*Lambda).DoCapture(id.Name, v)
		id.Var = usingScope.Find(id.Name)
	}
}

func (n *Lambda) Find(name string) Var {
	for _, a := range n.Args {
		assert(a.Var != nil)
		if name == a.Name {
			return a
		}
	}
	for _, a := range n.CapArgs {
		assert(a.Var != nil)
		if name == a.Name {
			return a
		}
	}
	return nil
}

func (n *Lambda) DoCapture(name string, v Var) {
	index := len(n.Args) + len(n.CapArgs)
	n.CapArgs = append(n.Args, &Ident{name, Local{index}})
	n.Captured = append(n.Captured, v)
}

type Frame interface {
	Find(name string) Var
}

type Frames []Frame

func (s *Frames) Push(v Frame) {
	*s = append(*s, v)
}

func (s *Frames) Pop() Frame {
	v := s.Last()
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *Frames) Last() Frame {
	return (*s)[len(*s)-1]
}

func (f *Frames) Find(name string) (Var, Frame) {
	s := *f
	for i := len(s) - 1; i >= 0; i-- {
		s := s[i]
		if v := s.Find(name); v != nil {
			return v, s
		}
	}
	return nil, nil
}

// gather records in s all definitions
// in the AST rooted at n.
//func gather(n Node) {
//	switch n := n.(type) {
//	default:
//		panic(unhandled(n))
//	case *Num, *Ident, *Call:
//		for _, n := range n.Args {
//			gather(n)
//		}
//	case *Lambda:
//		for _, id := range n.Args {
//			s.Def(id.Name, id)
//		}
//		gatherDefs(s, n.Body)
//	}
//}

//func resolve(s *Scope, n Node) {
//	switch n := n.(type) {
//	default:
//		panic(unhandled(n))
//	case *Num:
//	case *Ident:
//		if def := s.Resolve(n.Name); def != nil {
//			n.ID = def.ID
//		}
//	case *Call:
//		resolve(s, n.F)
//		for _, n := range n.Args {
//			resolve(s, n)
//		}
//	case *Lambda:
//		s = n.scope
//		resolve(s, n.Body)
//	}
//}
