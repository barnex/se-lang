package eva

type Const struct {
	v Value
}

func (c Const) Eval(m *Machine) {
	m.RA = c.v
}
