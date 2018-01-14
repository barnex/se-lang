package se

type Const struct {
	v Value
}

var _ Prog = (*Const)(nil)

func (c Const) Eval() Value {
	return c.v
}
