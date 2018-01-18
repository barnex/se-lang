package eva

import (
	"fmt"

	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/ast"
)

func compileIdent(id *ast.Ident) Prog {
	switch n := id.Var.(type) {
	default:
		panic(unhandled(n))
	case nil:
		return compileGlobal(id)
	case *ast.LocalVar:
		return compileLocalVar(n)
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

func compileLocalVar(n *ast.LocalVar) Prog {
	return &FromEBP{-2 - n.Index}
}

func compileLocal(i int) Prog {
	return &FromEBP{-2 - i}
}

type FromEBP struct {
	Offset int
}

func (p *FromEBP) Eval(s *Machine) {
	s.RA = s.FromBP(p.Offset, "local")
	fmt.Println("eval local", p.Offset, "RA=", s.RA)
}
