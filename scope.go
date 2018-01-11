package se

type Scope struct {
	parent  *Scope
	symbols map[string]*Ident
}

var identid = 0

func (e *Scope) Resolve(name string) *Ident {
	if n, ok := e.symbols[name]; ok {
		return n
	}
	if e.parent == nil {
		//panic(SyntaxErrorf("undefined: %v", name))
		return nil
	}
	return e.parent.Resolve(name)
}

func (e *Scope) Def(name string, value *Ident) {
	if _, ok := e.symbols[name]; ok {
		panic(SyntaxErrorf("already defined: %v", name))
	}
	identid++
	value.ID = identid
	e.symbols[name] = value
}

func (s *Scope) New() *Scope {
	return &Scope{parent: s, symbols: make(map[string]*Ident)}
}
