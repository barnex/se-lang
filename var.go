package se

type Var interface {
}

type Local struct {
	Index int
}

var _ Var = &Local{}

type Global struct {
	Name string // TODO
}
