package se

import (
	"fmt"
	"io"
)

type Stack struct {
	ID    int
	stack []Value
}

var stackID int

func NewStack() *Stack {
	stackID++
	return &Stack{stackID, nil}
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

func (s *Stack) PrintTo(w io.Writer) {
	fmt.Fprint(w, ":", s.ID)
}

func (s *Stack) AddStack(delta int) {
	new := len(s.stack) + delta
	if new > cap(s.stack) {
		s.stack = append(s.stack, make([]Value, new-cap(s.stack))...)
	}
	s.stack = s.stack[:new]
}
