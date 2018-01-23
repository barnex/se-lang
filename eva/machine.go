package eva

import "fmt"

type Machine struct {
	s  []Value
	ra Value
	bp int
}

func (m *Machine) SP() int {
	return len(m.s)
}

func (m *Machine) Push(v Value) {
	Log("push", v)
	m.s = append(m.s, v)
}

func (m *Machine) Pop() Value {
	v := m.s[len(m.s)-1]
	Log("pop", v)
	m.s = m.s[:len(m.s)-1]
	return v
}

func (m *Machine) FromBP(delta int) Value {
	v := m.s[m.bp+delta]
	Log("fromBP", delta, "=", v)
	return v
}

func (m *Machine) SetFromBP(delta int, v Value) {
	m.s[m.bp+delta] = v
	Log("setfromBP", delta, "to", v)
}

func (m *Machine) FromSP(delta int) Value {
	v := m.s[m.SP()+delta]
	Log("fromSP", delta, "=", v)
	return v
}

func (m *Machine) SetRA(v Value) {
	Log("ra=", v)
	m.ra = v
}

func (m *Machine) RA() Value {
	return m.ra
}

func (m *Machine) BP() int {
	return m.bp
}

func (m *Machine) SetBP(bp int) {
	m.bp = bp
}

func (m *Machine) Grow(delta int) {
	Log("grow", delta)
	new := len(m.s) + delta
	if new >= cap(m.s) {
		m.s = append(m.s, make([]Value, 1+new-cap(m.s))...)
	}
	m.s = m.s[:new]
}

func Log(x ...interface{}) {
	fmt.Println(x...)
}
