package se

type Var interface {
	Get() Value // TODO
}

var (
	_ Var = &Local{}
	_ Var = &Global{}
)

type Local struct {
	Index int
}

func (l *Local) Get() Value { panic("todo") }

type Global struct {
	Name string // TODO
}

func (l *Global) Get() Value { panic("todo") }
