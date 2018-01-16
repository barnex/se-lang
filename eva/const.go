package eva

type Const struct {
	v Value
}

func (c Const) Eval(s *Stack) {
	s.Push(c.v)
}
