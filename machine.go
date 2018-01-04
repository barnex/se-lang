package e

type Machine struct {
	stack []Node
}

func (m *Machine) Push(n Node) {
	m.stack = append(m.stack, n)
}

func (m *Machine) Pop() Node {
	n := m.Get(-1)
	m.AddStack(-1)
	return n
}

func (m *Machine) Get(off int) Node {
	return m.stack[len(m.stack)-off]
}

func (m *Machine) AddStack(delta int) {
	new := len(m.stack) + delta
	if new > cap(m.stack) {
		m.stack = append(m.stack, make([]Node, new-cap(m.stack))...)
	}
	m.stack = m.stack[:new]
}
