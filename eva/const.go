package eva

type Const struct {
	v Value
}

func (c Const) Eval(s *Machine) {
	s.Push(c.v, "const")
}
