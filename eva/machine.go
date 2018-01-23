package eva

import "fmt"

type Box struct {
	v *Value
}

func (b Box) Get() Value {
	return *b.v
}

func (b Box) Set(v Value) {
	if _, ok := v.(Box); ok {
		panic("boxing box")
	}
	*b.v = v
}

func box(v Value) Box {
	if v == nil {
		panic("boxing nil")
	}
	if _, ok := v.(Box); ok {
		panic("boxing box")
	}
	return Box{&v}
}

func (b Box) String() string {
	if b.v == nil {
		return "box(nil!)"
	} else {
		return fmt.Sprintf("box(%#v)", *b.v)
	}
}

type Machine struct {
	s  []Box
	ra Box
	bp int
}

func (m *Machine) SP() int {
	return len(m.s)
}

func (m *Machine) Push(b Box) {
	Log("push", b)
	m.s = append(m.s, b)
}

func (m *Machine) Pop() Box {
	v := m.s[len(m.s)-1]
	Log("pop", v)
	m.s = m.s[:len(m.s)-1]
	return v
}

func (m *Machine) FromBP(delta int) Box {
	v := m.s[m.bp+delta]
	Log("fromBP", delta, "=", v)
	return v
}

//func (m *Machine) SetFromBP(delta int, v Value) {
//	m.s[m.bp+delta] = v
//	Log("setfromBP", delta, "to", v)
//}

func (m *Machine) FromSP(delta int) Box {
	v := m.s[m.SP()+delta]
	Log("fromSP", delta, "=", v)
	return v
}

func (m *Machine) SetRA(v Box) {
	Log("ra=", v)
	m.ra = v
}

func (m *Machine) RA() Box {
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
	newl := len(m.s) + delta
	//if new >= cap(m.s) {
	//	m.s = append(m.s, make([]Box, 1+new-cap(m.s))...)
	//}
	for i := len(m.s) - 1; i < newl; i++ {
		m.s = append(m.s, Box{new(Value)})
	}
	m.s = m.s[:newl] // in case we shrink
}

func Log(x ...interface{}) {
	//fmt.Println(x...)
}
