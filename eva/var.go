package eva

import (
	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/ast"
)

var prelude = map[string]Prog{
	"add": fn(add),
}

func compileIdent(id *ast.Ident) Prog {
	switch n := id.Var.(type) {
	default:
		panic(unhandled(n))
	case nil:
		panic(se.Errorf("undefined: %q", id.Name))
	case *ast.LocalVar:
		return compileLocal(n)
	case *ast.GlobVar:
		return compileGlobal(n)
	}
}

func compileGlobal(n *ast.GlobVar) Prog {
	v, ok := prelude[n.Name]
	if !ok {
		panic(se.Errorf("undefined: %q", n.Name))
	}
	return v
}

func compileLocal(n *ast.LocalVar) Prog {
	return &FromTop{n.Index}
}

type FromTop struct {
	Index int
}

func (p *FromTop) Eval(s *Stack) {
	s.Push(s.FromTop(p.Index))
}

func add(s *Stack) {
	a := s.Pop().(float64)
	b := s.Pop().(float64)
	s.Push(a + b)
}

type fn func(*Stack)

func (f fn) Eval(s *Stack)  { s.Push(f) }
func (f fn) Apply(s *Stack) { f(s) }
