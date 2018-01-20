package eva

var prelude = map[string]Prog{
	"add": fn(add),
	"mul": fn(mul),
}

func add(m *Machine) {
	a := m.FromBP(-2, "a").(float64)
	b := m.FromBP(-3, "b").(float64)
	m.RA = a + b
}

func mul(m *Machine) {
	a := m.FromBP(-2, "a").(float64)
	b := m.FromBP(-3, "b").(float64)
	m.RA = a * b
}

type fn func(*Machine)

func (f fn) Exec(m *Machine)  { m.RA = f }
func (f fn) Apply(m *Machine) { f(m) }
func (f fn) NFrame() int      { return 2 }

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}

//var _ Applier = fn(nil)
