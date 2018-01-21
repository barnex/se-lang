package eva

var prelude = pkg{
	"add":   fn2(add),
	"and":   fn2(and),
	"eq":    fn2(eq),
	"false": &Const{false},
	"ge":    fn2(ge),
	"gt":    fn2(gt),
	"le":    fn2(le),
	"lt":    fn2(lt),
	"mul":   fn2(mul),
	"neq":   fn2(neq),
	"or":    fn2(or),
	"true":  &Const{true},
}

type pkg map[string]Prog

func (p pkg) Find(name string) Prog {
	return p[name]
}

type fn2 func(a, b Value) Value

func (f fn2) Exec(m *Machine) {
	m.SetRA(f)
}

func (f fn2) Apply(m *Machine) {
	a := m.FromSP(-1)
	b := m.FromSP(-2)
	m.SetRA(f(a, b))
}

func add(a, b Value) Value { return a.(float64) + b.(float64) }
func and(a, b Value) Value { return a.(bool) && b.(bool) }
func eq(a, b Value) Value  { return a == b }
func ge(a, b Value) Value  { return a.(float64) >= b.(float64) }
func gt(a, b Value) Value  { return a.(float64) > b.(float64) }
func le(a, b Value) Value  { return a.(float64) <= b.(float64) }
func lt(a, b Value) Value  { return a.(float64) < b.(float64) }
func mul(a, b Value) Value { return a.(float64) * b.(float64) }
func neq(a, b Value) Value { return a != b }
func or(a, b Value) Value  { return a.(bool) || b.(bool) }
