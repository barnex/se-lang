package eva

var prelude = pkg{
	"add":   fn2(add),
	"sub":   fn2(sub),
	"and":   fn2(and),
	"eq":    fn2(eq),
	"false": &Const{false},
	"ge":    fn2(ge),
	"gt":    fn2(gt),
	"le":    fn2(le),
	"lt":    fn2(lt),
	"mod":   fn2(mod),
	"mul":   fn2(mul),
	"neg":   fn1(neg),
	"neq":   fn2(neq),
	"not":   fn1(not),
	"or":    fn2(or),
	"true":  &Const{true},
}

type pkg map[string]Prog

func (p pkg) Find(name string) Prog {
	return p[name]
}

type fn1 func(a Value) Value

func (f fn1) Exec(m *Machine) {
	m.SetRA(box(f))
}

func (f fn1) Apply(m *Machine) {
	a := m.FromSP(-1).Get()
	m.SetRA(box(f(a)))
}

func neg(a Value) Value { return -a.(int) }
func not(a Value) Value { return !a.(bool) }

type fn2 func(a, b Value) Value

func (f fn2) Exec(m *Machine) {
	m.SetRA(box(f))
}

func (f fn2) Apply(m *Machine) {
	a := m.FromSP(-1).Get()
	b := m.FromSP(-2).Get()
	m.SetRA(box(f(a, b)))
}

func add(a, b Value) Value { return a.(int) + b.(int) }
func and(a, b Value) Value { return a.(bool) && b.(bool) }
func eq(a, b Value) Value  { return a == b }
func ge(a, b Value) Value  { return a.(int) >= b.(int) }
func gt(a, b Value) Value  { return a.(int) > b.(int) }
func le(a, b Value) Value  { return a.(int) <= b.(int) }
func lt(a, b Value) Value  { return a.(int) < b.(int) }
func mul(a, b Value) Value { return a.(int) * b.(int) }
func neq(a, b Value) Value { return a != b }
func or(a, b Value) Value  { return a.(bool) || b.(bool) }
func sub(a, b Value) Value { return a.(int) - b.(int) }
func mod(a, b Value) Value { return a.(int) % b.(int) }
