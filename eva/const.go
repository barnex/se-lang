package eva

type Const struct {
	v Value
}

func (c Const) Eval() Value {
	return c.v
}
