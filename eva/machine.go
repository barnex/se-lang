package eva

import "fmt"

type Machine struct {
	s  []Value
	RA Value
	BP int
}

func (m *Machine) SP() int {
	return len(m.s)
}

func (m *Machine) Push(v Value) {
	fmt.Println("push", v)
	m.s = append(m.s, v)
}

func (m *Machine) Pop(msg string) Value {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("<-pop :", msg, ":", err)
		}
	}()
	v := m.s[len(m.s)-1]
	fmt.Println("<-pop :", v, msg, "(len=", len(m.s), ")")
	m.s = m.s[:len(m.s)-1]
	return v
}

func (m *Machine) FromBP(delta int) Value {
	v := m.s[m.BP+delta]
	fmt.Printf("fromBP %v=%v \n", delta, v)
	return v
}

func (m *Machine) FromSP(delta int) Value {
	v := m.s[m.SP()+delta]
	fmt.Printf("fromSP %v=%v \n", delta, v)
	return v
}

func (m *Machine) Grow(delta int) {
	fmt.Println("grow", delta)
	new := len(m.s) + delta
	if new >= cap(m.s) {
		m.s = append(m.s, make([]Value, 1+new-cap(m.s))...)
	}
	m.s = m.s[:new]
}
