package eva

import (
	"fmt"

	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/ast"
)

// Var refers to the storage location of a variable,
// so we can set or retreive the its value.
type Var interface {
	Prog
}

type Arg struct {
	Index int
}

var _ Var = (*Arg)(nil)

func (a *Arg) Exec(m *Machine) {
	m.RA = m.FromBP(-2 - a.Index)
}

func (a *Arg) String() string {
	return fmt.Sprint("$", a.Index)
}

// A CaptVar refers to a captured variable:
// a variable closed over by a closure.
//type CaptVar struct {
//	Name string
//	Src  Var // variable being captured from the parent frame
//	Dst  Var // variable being captured to
//}

//func (c *CaptVar) String() string { return fmt.Sprintf("[%v=p.%v]", c.Name, c.Src) }

// A LocalVar refers to a local variable:
// a variable that exist on a call stack (argument or local define)
//type LocVar struct {
//	Name  string
//	Index int
//}

//func (l *LocVar) variable()      {}
//func (l *LocVar) String() string { return fmt.Sprint("L", l.Index) }

func compileIdent(id *ast.Ident) Prog {
	switch n := id.Object.(type) {
	default:
		panic(unhandled(n))
	case nil:
		panic(se.Errorf("compileIdent: undefined: %q", id.Name))
	case Prog:
		return n
	}
}

func compileArg(a *Arg) Prog {
	panic("todo")
}

//func compileGlobal(id *ast.Ident) Prog {
//	assert(id.Var == nil)
//	v, ok := prelude[id.Name]
//	if !ok {
//		panic(se.Errorf("undefined: %q", id.Name))
//	}
//	return v
//}
//
//func compileLocalVar(n *ast.LocalVar) Prog {
//	return &FromEBP{-2 - n.Index}
//}
//
//func compileLocal(i int) Prog {
//	return &FromEBP{-2 - i}
//}
//
//type FromEBP struct {
//	Offset int
//}
//
//func (p *FromEBP) Eval(s *Machine) {
//	s.RA = s.FromBP(p.Offset, "local")
//	fmt.Println("eval local", p.Offset, "RA=", s.RA)
//}
//
//type Var interface {
//	variable()
//}
