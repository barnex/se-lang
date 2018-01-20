package eva

import (
	"log"

	"github.com/barnex/se-lang/ast"
)

func Resolve(n ast.Node) {
	resolve(Frames{prelude}, n)
}

func resolve(s Frames, n ast.Node) {
	Log("resolve", n)
	switch n := n.(type) {
	case *ast.Call:
		resolveCall(s, n)
	case *ast.Ident:
		resolveIdent(s, n)
	case *ast.Lambda:
		resolveLambda(s, n)
	case *ast.Num:
		// nothing to do
	default:
		panic(unhandled(n))
	}
}

func resolveCall(s Frames, c *ast.Call) {
	resolve(s, c.F)
	for _, a := range c.Args {
		resolve(s, a)
	}
}

func resolveIdent(s Frames, id *ast.Ident) {
	Log("resolveIdent", id)

	// defScope: where ident was defined
	// usingScope: where ident is being used:
	// 	x ->           // defScope of x
	// 		y ->
	//  		x + y  // usingScope of x

	name := id.Name
	v, defScope := s.Find(name)
	if v == nil {
		Log("resolveIdent: not found", id)
		return
	}

	switch {
	case defScope == -1: // not found
		// leave open for now, compile will search for global
	case defScope == len(s)-1: // directly under parent
		id.Object = v
		//id.Parent = s[defScope]
		//default: // captured variable
		//	// loop over frames, capture from defscope+1 to last, capture all the way
		//	for i := defScope + 1; i < len(s); i++ {
		//		v := s[i-1].Find(name)
		//		//s[i].(*Lambda).DoCapture(name, v)
		//	}
		//	v := s[len(s)-1].Find(name)
		//	id.Object = v
		//	//id.Parent = s[defScope]
	}
}

func resolveLambda(s Frames, n *ast.Lambda) {
	// first define the arguments
	for i, a := range n.Args {
		a.Object = &Arg{Index: i}
	}

	// then resolve the body
	s.Push(&LambdaFrame{n.Args})
	resolve(s, n.Body)
	s.Pop()
}

type LambdaFrame struct {
	Args []*ast.Ident
}

func (n *LambdaFrame) Find(name string) Var {
	for _, a := range n.Args {
		assert(a.Object != nil)
		if name == a.Name {
			return a.Object.(Var)
		}
	}
	//for _, a := range n.Caps {
	//	if name == a.Name {
	//		return a
	//	}
	//}
	return nil // not found, maybe global
}

// TODO: should not be method
//func (n *Lambda) DoCapture(name string, v Var) {
//	if v := n.Find(name); v != nil {
//		return // already captured
//	}
//	c := &CaptVar{
//		Name: name,
//		Src:  v,
//		//Dst:  &CaptVar{},
//	}
//	n.Caps = append(n.Caps, c)
//}

//func (n *Lambda) NumLocals() int {
//	return len(n.Args) + len(n.Caps)
//}

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

func (f *Frames) Find(name string) (Var, int) {
	s := *f
	for i := len(s) - 1; i >= 0; i-- {
		if v := s[i].Find(name); v != nil {
			return v, i
		}
	}
	return nil, 0
}

func Log(action string, arg interface{}) {
	log.SetFlags(0)
	log.Printf("%s: %#v\n", action, arg)
}
