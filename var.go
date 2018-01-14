package se

// Var refers to the storage location of a variable,
// so we can set or retreive the its value.
type Var interface {
	Get() Value
}

var (
	_ Var = (*CaptVar)(nil)
	_ Var = (*GlobVar)(nil)
	_ Var = (*LocalVar)(nil)
)

// A CaptVar refers to a captured variable:
// a variable closed over by a closure.
type CaptVar struct {
	Name      string
	ParVar    *LocalVar // variable being captured from the parent frame
	*LocalVar           // local variable being captured to
}

// A GlobVar refers to a global variabe:
// a variable with a constant address.
type GlobVar struct {
	Name string
}

func (l *GlobVar) Get() Value {
	panic("todo")
}

// A LocalVar refers to a local variable:
// a variable that exist on a call stack (argument or local define)
type LocalVar struct {
	Index int
}

func (l *LocalVar) Get() Value {
	panic("todo")
}
