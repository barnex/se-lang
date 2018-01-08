package se

type Stack struct {
	stack []Value
}

func (s *Stack) Push(v Value) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() Value {
	v := s.stack[len(s.stack)-1]
	s.AddStack(-1)
	return v
}

func (s *Stack) Eval() Value {
	return s.stack[len(s.stack)-1]
}

//func (s *Stack) Get(off int) Node {
//	return m.stack[len(m.stack)-off]
//}

func (s *Stack) AddStack(delta int) {
	new := len(s.stack) + delta
	if new > cap(s.stack) {
		s.stack = append(s.stack, make([]Value, new-cap(s.stack))...)
	}
	s.stack = s.stack[:new]
}
