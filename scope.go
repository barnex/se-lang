package se

type Scope struct {
	parent  *Scope
	symbols map[string]Obj
}

type Obj interface{}

func (e *Scope) Resolve(name string) Obj {
	if n, ok := e.symbols[name]; ok {
		return n
	}
	if e.parent == nil {
		panic(SyntaxErrorf("undefined: %v", name))
	}
	return e.parent.Resolve(name)
}

func (e *Scope) Def(name string, value Obj) {
	if _, ok := e.symbols[name]; ok {
		panic(SyntaxErrorf("already defined: %v", name))
	}
	e.symbols[name] = value
}

func (s *Scope) New() *Scope {
	return &Scope{parent: s, symbols: make(map[string]Obj)}
}
