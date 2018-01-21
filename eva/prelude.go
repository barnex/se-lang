package eva

var prelude = pkg{
	"add": fn(add),
	"mul": fn(mul),
}

type pkg map[string]Prog

func (p pkg) Find(name string) Prog {
	return p[name]
}

func add(m *Machine) {
	a := m.FromSP(-1).(float64)
	b := m.FromSP(-2).(float64)
	m.SetRA(a + b)
}

func mul(m *Machine) {
	a := m.FromSP(-1).(float64)
	b := m.FromSP(-2).(float64)
	m.SetRA(a * b)
}

type fn func(*Machine)

func (f fn) Exec(m *Machine)  { m.SetRA(f) }
func (f fn) Apply(m *Machine) { f(m) }
func (f fn) NFrame() int      { return 2 }

//var _ Applier = fn(nil)
