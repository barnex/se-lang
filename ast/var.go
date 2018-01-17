package ast

import "fmt"

// Var refers to the storage location of a variable,
// so we can set or retreive the its value.
type Var interface {
	variable()
}

var (
	_ Var = (*CaptVar)(nil)
	_ Var = (*GlobVar)(nil)
	_ Var = (*LocalVar)(nil)
)

// A CaptVar refers to a captured variable:
// a variable closed over by a closure.
type CaptVar struct {
	Name   string
	ParVar *LocalVar // variable being captured from the parent frame
	Local  *LocalVar // local variable being captured to
}

func (c *CaptVar) variable()      {}
func (c *CaptVar) String() string { return fmt.Sprintf("[%v=p.%v]", c.Local, c.ParVar) }

// A GlobVar refers to a global variabe:
// a variable with a constant address.
type GlobVar struct {
	Name string
}

func (l *GlobVar) variable()      {}
func (l *GlobVar) String() string { return "$$" }

// A LocalVar refers to a local variable:
// a variable that exist on a call stack (argument or local define)
type LocalVar struct {
	Index int
}

func (l *LocalVar) variable()      {}
func (l *LocalVar) String() string { return fmt.Sprint("$", l.Index) }
