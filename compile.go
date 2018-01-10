package se

import (
	"fmt"
	"io"
)

func Compile(src io.Reader) (_ Prog, e error) {
	ast, err := Parse(src)
	if err != nil {
		return nil, err
	}

	defer func() {
		switch p := recover().(type) {
		default:
			panic(p)
		case nil:
		case *SyntaxError:
			e = p
		}
	}()

	return compileProg(ast), nil
}

type Prog interface {
	Eval() interface{}
	Node
}

func compileProg(n Node) Prog {
	toplevel := prelude.New()
	gatherDefs(toplevel, n)
	return compile(toplevel, n)
}

// gatherDefs records in s all definitions
// in the AST rooted at n.
func gatherDefs(s *Scope, n Node) {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *Num, *Ident: // nothing to do
	case *Call:
		gatherDefs(s, n.F)
		for _, n := range n.Args {
			gatherDefs(s, n)
		}
	case *Lambda:
		s = s.New()
		n.scope = s
		for _, id := range n.Args {
			s.Def(id.Name, NewStack())
		}
		gatherDefs(s, n.Body)
	}
}

func compile(s *Scope, n Node) Prog {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *Num:
		return n
	case *Ident:
		return compileIdent(s, n)
	case *Call:
		return compileCall(s, n)
	case *Lambda:
		return compileLambda(n.scope, n)
	}
}

func compileCall(s *Scope, n *Call) Prog {
	args := make([]Prog, len(n.Args))
	for i, a := range n.Args {
		args[i] = compile(s, a)
	}
	f := compile(s, n.F) // todo message
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

func (n *PCall) PrintTo(w io.Writer) {
	n.F.PrintTo(w)
	fmt.Fprint(w, "(")
	for i, a := range n.Args {
		if i != 0 {
			fmt.Fprint(w, ", ")
		}
		a.PrintTo(w)
	}
	fmt.Fprint(w, ")")
}

func compileLambda(s *Scope, n *Lambda) Prog {
	args := make([]*Stack, len(n.Args))
	for i := range args {
		args[i] = s.Resolve(n.Args[i].Name).(*Stack)
	}
	return &PLambda{
		Args: args,
		Body: compile(n.scope, n.Body),
	}
}

type PLambda struct {
	Args []*Stack
	Body Prog
}

func (n *PLambda) Eval() interface{} {
	return n
}

func (n *PLambda) PrintTo(w io.Writer) {
	fmt.Fprint(w, "(")
	for i, a := range n.Args {
		if i != 0 {
			fmt.Fprint(w, ", ")
		}
		a.PrintTo(w)
	}
	fmt.Fprint(w, ")->(")
	n.Body.PrintTo(w)
	fmt.Fprint(w, ")")
}

func (n *PLambda) Apply(args []Value) Value {
	for i, a := range n.Args {
		a.Push(args[i])
	}
	v := n.Body.Eval()
	for _, a := range n.Args {
		a.Pop()
	}
	return v
}

func compileIdent(s *Scope, n *Ident) Prog {
	return s.Resolve(n.Name).(Prog)
}

type Applier interface {
	Apply([]Value) Value
}

func unhandled(n Node) string {
	return fmt.Sprintf("bug: unhandled AST node: %T", n)
}
