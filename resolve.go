package se

import "fmt"

func Resolve(n Node) {
	scope := prelude.New()
	gatherDefs(scope, n)
	resolve(scope, n)
}

// gatherDefs records in s all definitions
// in the AST rooted at n.
func gatherDefs(s *Scope, n Node) {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *Num, *Ident: // nothing to do
	case *Call:
		gatherDefs(s, n.F)
		for _, n := range n.Args {
			gatherDefs(s, n)
		}
	case *Lambda:
		s = s.New()
		n.scope = s
		for _, id := range n.Args {
			s.Def(id.Name, id)
		}
		gatherDefs(s, n.Body)
	}
}

func resolve(s *Scope, n Node) {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *Num:
	case *Ident:
		if def := s.Resolve(n.Name); def != nil {
			n.ID = def.ID
		}
	case *Call:
		resolve(s, n.F)
		for _, n := range n.Args {
			resolve(s, n)
		}
	case *Lambda:
		s = n.scope
		resolve(s, n.Body)
	}
}

func unhandled(n Node) string {
	return fmt.Sprintf("bug: unhandled AST node: %T", n)
}
