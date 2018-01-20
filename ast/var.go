package ast

import "fmt"

// Var refers to the storage location of a variable,
// so we can set or retreive the its value.
type Var interface {
	variable()
}

var (
	_ Var = (*Arg)(nil)
	_ Var = (*CaptVar)(nil)
)

// A CaptVar refers to a captured variable:
// a variable closed over by a closure.
type CaptVar struct {
	Name string
	Src  Var // variable being captured from the parent frame
	Dst  Var // variable being captured to
}

func (c *CaptVar) variable()      {}
func (c *CaptVar) String() string { return fmt.Sprintf("[%v=p.%v]", c.Name, c.Src) }

type Arg struct {
	Name  string
	Index int
}

func (l *Arg) variable()      {}
func (l *Arg) String() string { return fmt.Sprint("$", l.Index) }

// A LocalVar refers to a local variable:
// a variable that exist on a call stack (argument or local define)
//type LocVar struct {
//	Name  string
//	Index int
//}

//func (l *LocVar) variable()      {}
//func (l *LocVar) String() string { return fmt.Sprint("L", l.Index) }
