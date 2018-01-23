package ast

import (
	"fmt"
)

// Resolve traverses the AST and populates the Var fields of all identifiers.
// Lambda arguments are resolved to *Arg.
// Local variables are resolved to *LocVar.
func Resolve(n Node) {
	gather(n, Frames{})
	resolve(Frames{}, n)
}

// gather traverses the AST and records all variable declarations.
// The declaring identifier's Var field is set to a new variable. E.g.:
// 	x -> x + 1
//  |
// 	\
// 	 Arg{0}
func gather(n Node, s Frames) {
	switch n := n.(type) {
	case *Assign:
		gatherAssign(n, s)
	case *Block:
		gatherBlock(n, s)
	case *Call:
		gatherCall(n, s)
	case *Cond:
		gatherCond(n, s)
	case *Ident:
		gatherIdent(n, s)
	case *Lambda:
		gatherLambda(n, s)
	case *Num: // nothing to do
	default:
		panic(unhandled(n))
	}
}

// resolve traverses the AST and resolves all identifiers that are not declarations.
// It sets their Var field to that of the corresponding declaration.
// gather must have been called first.
// 	E.g.:
// 	x -> x + 1
//  |    |
// 	\   /
// 	 Arg{0}
func resolve(s Frames, n Node) {
	switch n := n.(type) {
	case *Assign:
		resolveAssign(s, n)
	case *Block:
		resolveBlock(s, n)
	case *Call:
		resolveCall(s, n)
	case *Cond:
		resolveCond(s, n)
	case *Ident:
		resolveIdent(s, n)
	case *Lambda:
		resolveLambda(s, n)
	case *Num: // nothing to do
	default:
		panic(unhandled(n))
	}
}

// ---- Assign

func gatherAssign(a *Assign, s Frames) {
	a.LHS.Var = parentLambda(s).NewVariable()
	gather(a.RHS, s)
}

func resolveAssign(s Frames, a *Assign) {
	resolve(s, a.RHS)
}

// ---- Block

func gatherBlock(b *Block, s Frames) {
	s.Push(b)
	defer s.Pop()

	for _, stmt := range b.Stmts {
		gather(stmt, s)
	}
}

func resolveBlock(s Frames, b *Block) {
	s.Push(b)
	defer s.Pop()

	for _, stmt := range b.Stmts {
		resolve(s, stmt)
	}
}

func (b *Block) Find(name string) Var {
	for _, stmt := range b.Stmts {
		if a, ok := stmt.(*Assign); ok {
			if a.LHS.Name == name {
				assert(a.LHS.Var != nil)
				return a.LHS.Var
			}
		}
	}
	return nil
}

// ---- Call

func gatherCall(c *Call, s Frames) {
	gather(c.F, s)
	for _, a := range c.Args {
		gather(a, s)
	}
}

func resolveCall(s Frames, c *Call) {
	resolve(s, c.F)
	for _, a := range c.Args {
		resolve(s, a)
	}
}

// ---- Cond

func gatherCond(n *Cond, s Frames) {
	gather(n.Test, s)
	gather(n.If, s)
	gather(n.Else, s)
}

func resolveCond(s Frames, n *Cond) {
	resolve(s, n.Test)
	resolve(s, n.If)
	resolve(s, n.Else)
}

// ---- Ident
// Here be dragons.

func gatherIdent(id *Ident, s Frames) {

	name := id.Name
	v, defScope := s.Find(name)
	if v == nil {
		return // not in parent, so not a capture
	}

	switch {
	case defScope == -1:
		// not found
	case defScope == len(s)-1:
		// argument
	default:
		// captured variable
		// loop over frames, capture from defscope+1 to last, capture all the way
		for i := defScope + 1; i < len(s); i++ {
			if l, ok := s[i].(*Lambda); ok {
				v := s[i-1].Find(name)
				l.DoCapture(name, v)
			}
		}
	}
}

func resolveIdent(s Frames, id *Ident) {
	v, _ := s.Find(id.Name)
	id.Var = v
}

// ---- Lambda

func gatherLambda(n *Lambda, s Frames) {
	s.Push(n)
	defer s.Pop()

	for i, a := range n.Args {
		a.Var = &Arg{Index: i}
	}
	gather(n.Body, s)
}

func resolveLambda(s Frames, n *Lambda) {
	s.Push(n)
	defer s.Pop()

	resolve(s, n.Body)
}

func parentLambda(s Frames) *Lambda {
	//if len(s) < 2 {
	//	panic("no parent frame (1)")
	//}
	for i := len(s) - 1; i >= 0; i-- {
		if l, ok := s[i].(*Lambda); ok {
			return l
		}
	}
	panic("no parent frame (2)")
}

func (n *Lambda) Find(name string) Var {
	Log("lambdaframe: find", name)
	for _, a := range n.Args {
		if name == a.Name {
			Log("lambdaframe: found", a.Var)
			assert(a.Var != nil)
			return a.Var
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

func (n *Lambda) DoCapture(name string, v Var) {
	assert(v != nil)
	if v := n.Find(name); v != nil {
		return // already captured
	}
	c := Capture{
		Name: name,
		Src:  v,
		Dst:  n.NewVariable(),
	}
	n.Caps = append(n.Caps, c)
	Log("lambdaframe: docapture:", c)
}

// ---------

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
