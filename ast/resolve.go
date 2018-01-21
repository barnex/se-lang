package ast

import (
	"fmt"
)

func Resolve(n Node) {
	resolve(Frames{}, n)
}

func resolve(s Frames, n Node) {
	Log("resolve", n)
	switch n := n.(type) {
	case *Call:
		resolveCall(s, n)
	case *Ident:
		resolveIdent(s, n)
	case *Lambda:
		resolveLambda(s, n)
	case *Num:
		// nothing to do
	case *Block:
		resolveBlock(s, n)
	default:
		panic(unhandled(n))
	}
}

func resolveBlock(s Frames, b *Block) {
	panic("todo")
	//s.Push(b)
	//defer s.Pop()

	//for _, stmt := range b.Stmts {
	//	if a, ok := stmt.(*Assign); ok {
	//		resolve(s, a.RHS)
	//	} else {
	//		resolve(s, stmt)
	//	}
	//}
}

//func (b*Block) Find() Var{
//
//}

func resolveCall(s Frames, c *Call) {
	Log("resolveCall", c)
	resolve(s, c.F)
	for _, a := range c.Args {
		resolve(s, a)
	}
}

func resolveIdent(s Frames, id *Ident) {
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
		id.Var = v
	//case defScope == 0: // global
	//	id.Var = v
	default: // captured variable
		// loop over frames, capture from defscope+1 to last, capture all the way
		for i := defScope + 1; i < len(s); i++ {
			v := s[i-1].Find(name)
			s[i].(LambdaFrame).DoCapture(name, v)
		}
		v := s[len(s)-1].Find(name)
		assert(v != nil)
		id.Var = v
	}
}

func resolveLambda(s Frames, n *Lambda) {
	Log("resolveLambda", n)
	// first define the arguments
	for i, a := range n.Args {
		a.Var = &Arg{Index: i}
	}

	// then resolve the body
	f := LambdaFrame{n}
	s.Push(f)
	defer s.Pop()
	resolve(s, n.Body)
}

type LambdaFrame struct {
	*Lambda
}

func (n LambdaFrame) Find(name string) Var {
	Log("lambdaframe: find", name)
	for _, a := range n.Args {
		if name == a.Name {
			Log("lambdaframe: found", a.Var)
			return a.Var.(Var)
		}
	}
	for _, c := range n.Caps {
		if name == c.Name {
			Log("lambdaframe: found: captured:", c.Dst)
			return c.Dst.(Var)
		}
	}

	Log("lambdaframe: not found", nil)
	return nil // not found, maybe global
}

func (n LambdaFrame) DoCapture(name string, v Var) {
	if v := n.Find(name); v != nil {
		return // already captured
	}
	c := Capture{
		Name: name,
		Src:  v,
		Dst:  &LocVar{len(n.Caps)},
	}
	n.Caps = append(n.Caps, c)
	Log("lambdaframe: docapture:", c)
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

func (f *Frames) Find(name string) (Var, int) {
	Log("frames: find", name)
	s := *f
	for i := len(s) - 1; i >= 0; i-- {
		if v := s[i].Find(name); v != nil {
			Log("frames: found", v)
			return v, i
		}
	}
	Log("frames: not found", nil)
	return nil, 0
}

func Log(action string, arg interface{}) {
	//log.SetFlags(0)
	//log.Printf("%s: %#v\n", action, arg)
}

func unhandled(x interface{}) string {
	return fmt.Sprintf("BUG: unhandled case: %T", x)
}

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}
