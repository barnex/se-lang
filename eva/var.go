package eva

import (
	"fmt"

	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/ast"
)

var prelude = map[string]Prog{
	"add": fn(add),
	"mul": fn(mul),
}

func compileIdent(id *ast.Ident) Prog {
	switch n := id.Var.(type) {
	default:
		panic(unhandled(n))
	case nil:
		return compileGlobal(id)
	case *ast.LocalVar:
		return compileLocal(n)
		//case *ast.GlobVar:
		//return compileGlobal(n)
	}
}

func compileGlobal(id *ast.Ident) Prog {
	assert(id.Var == nil)
	v, ok := prelude[id.Name]
	if !ok {
		panic(se.Errorf("undefined: %q", id.Name))
	}
	return v
}

func compileLocal(n *ast.LocalVar) Prog {
	return &FromEBP{-n.Index - 2}
}

type FromEBP struct {
	Offset int
}

func (p *FromEBP) Eval(s *Machine) {
	msg := fmt.Sprint("local ", p.Offset)
	s.Push(s.FromEBP(p.Offset, msg), msg)
}

func add(s *Machine) {
	a := s.FromEBP(-2, "a").(float64)
	b := s.FromEBP(-3, "b").(float64)
	s.Push(a+b, "a+b")
}

func mul(s *Machine) {
	a := s.FromEBP(-2, "a").(float64)
	b := s.FromEBP(-3, "b").(float64)
	s.Push(a*b, "a*b")
}

type fn func(*Machine)

func (f fn) Eval(s *Machine)  { s.Push(f, "func:self") }
func (f fn) Apply(s *Machine) { f(s) }

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}
