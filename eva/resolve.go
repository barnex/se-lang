package eva

import (
	"fmt"
	"log"

	"github.com/barnex/se-lang/ast"
)

func Resolve(n ast.Node) {
	resolve(Frames{prelude}, n)
	fmt.Println("*** resolved:", ast.ToString(n))
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
	Log("resolveCall", c)
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
	case defScope == 0: // global
		id.Object = v
	default: // captured variable
		// loop over frames, capture from defscope+1 to last, capture all the way
		for i := defScope + 1; i < len(s); i++ {
			v := s[i-1].Find(name)
			s[i].(*LambdaFrame).DoCapture(name, v)
		}
		v := s[len(s)-1].Find(name)
		id.Object = v
		//id.Parent = s[defScope]
	}
}

func resolveLambda(s Frames, n *ast.Lambda) {
	Log("resolveLambda", n)
	// first define the arguments
	for i, a := range n.Args {
		a.Object = &Arg{Index: i}
	}

	// then resolve the body
	f := &LambdaFrame{Args: n.Args}
	n.Object = f
	s.Push(f)
	resolve(s, n.Body)
	s.Pop()
}

type LambdaFrame struct {
	// TODO: just wrap lambda?
	Args []*ast.Ident
	Caps []Capture
	orig *ast.Lambda
}

type Capture struct {
	Name string
	Src  Var // variable being captured from the parent frame
	Dst  Var // variable being captured to
}

func (c Capture) String() string {
	return fmt.Sprint(c.Dst, "=", c.Src)
}

func (n *LambdaFrame) Find(name string) Var {
	Log("lambdaframe: find", name)
	for _, a := range n.Args {
		if name == a.Name {
			Log("lambdaframe: found", a.Object)
			return a.Object.(Var)
		}
	}
	for _, c := range n.Caps {
		if name == c.Name {
			Log("lambdaframe: found: captured:", c.Dst)
			return c.Dst
		}
	}

	Log("lambdaframe: not found", nil)
	return nil // not found, maybe global
}

func (n *LambdaFrame) DoCapture(name string, v Var) {
	if v := n.Find(name); v != nil {
		return // already captured
	}
	c := Capture{
		Src: v,
		Dst: &LocVar{len(n.Caps)},
	}
	n.Caps = append(n.Caps, c)
	fmt.Println("lambdaframe: docapture:", c)
}

func (n *LambdaFrame) String() string {
	str := "["
	for _, c := range n.Caps {
		str += c.String() + ","
	}
	str += "]"
	return str
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
	log.SetFlags(0)
	log.Printf("%s: %#v\n", action, arg)
}
