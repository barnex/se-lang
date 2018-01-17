package eva

import "fmt"

type Machine struct {
	s   []Value
	EBP int
}

func (m *Machine) ESP() int {
	return len(m.s)
}

func (m *Machine) Push(v Value, msg string) {
	fmt.Println("->push:", v, msg)
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

func (m *Machine) FromEBP(delta int, msg string) Value {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("fromebp %v: %v %v\n", delta, msg, err)
		}
	}()

	v := m.s[m.EBP+delta]
	fmt.Printf("fromebp %v: %v %v\n", delta, v, msg)
	return v
}

func (m *Machine) Grow(delta int) {
	new := len(m.s) + delta
	if new > cap(m.s) {
		m.s = append(m.s, make([]Value, new-cap(m.s))...)
	}
	m.s = m.s[:new]
}
