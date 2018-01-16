package eva

import "fmt"

type Stack []Value

func (s *Stack) Len() int {
	return len(*s)
}

func (s *Stack) Push(v Value) {
	fmt.Println("->push:", v)
	*s = append(*s, v)
}

func (s *Stack) Pop() Value {
	v := s.FromTop(-1)
	s.Grow(-1)
	fmt.Println("<-pop :", v)
	return v
}

func (s *Stack) FromTop(delta int) Value {
	return (*s)[s.Len()+delta]
}

func (s *Stack) Grow(delta int) {
	new := len(*s) + delta
	if new > cap(*s) {
		*s = append(*s, make([]Value, new-cap(*s))...)
	}
	*s = (*s)[:new]
}
