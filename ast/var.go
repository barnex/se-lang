package ast

import (
	"fmt"
)

// Var refers to the storage location of a variable,
// so we can set or retreive the its value.
type Var interface {
}

type Arg struct {
	Index int
}

func (a *Arg) String() string {
	return fmt.Sprint("$", a.Index)
}

type LocVar struct {
	Index int
}

func (l *LocVar) String() string {
	return fmt.Sprint("L", l.Index)
}

//func compileGlobal(id *Ident) Prog {
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

//func (a *Arg) Exec(m *Machine) {
//	m.SetRA(m.FromBP(-2 - a.Index))
//}

//func (l *LocVar) Exec(m *Machine) {
//	m.SetRA(m.FromBP(l.Index))
//}
