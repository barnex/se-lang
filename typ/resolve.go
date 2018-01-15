package ast

//func Resolve(n Node) {
//	resolve(Frames{prelude}, n)
//}
//
//func resolve(s Frames, n Node) {
//	switch n := n.(type) {
//	case *Ident:
//		resolveIdent(s, n)
//	case *Lambda:
//		resolveLambda(s, n)
//	case *Num, *Call:
//		for _, n := range children(n) {
//			resolve(s, n)
//		}
//	default:
//		panic(unhandled(n))
//	}
//}
//
//// Var refers to the storage location of a variable,
//// so we can set or retreive the its value.
//type Var interface {
//	variable()
//}
//
//var (
//	_ Var = (*CaptVar)(nil)
//	_ Var = (*GlobVar)(nil)
//	_ Var = (*LocalVar)(nil)
//)
//
//// A CaptVar refers to a captured variable:
//// a variable closed over by a closure.
//type CaptVar struct {
//	Name   string
//	ParVar *LocalVar // variable being captured from the parent frame
//	Local  *LocalVar // local variable being captured to
//}
//
//func (c *CaptVar) variable()      {}
//func (c *CaptVar) String() string { return fmt.Sprintf("[%v=p.%v]", c.Local, c.ParVar) }
//
//// A GlobVar refers to a global variabe:
//// a variable with a constant address.
//type GlobVar struct {
//	Name string
//}
//
//func (l *GlobVar) variable() {}
//
//// A LocalVar refers to a local variable:
//// a variable that exist on a call stack (argument or local define)
//type LocalVar struct {
//	Index int
//}
//
//func (l *LocalVar) variable()      {}
//func (l *LocalVar) String() string { return fmt.Sprint("$", l.Index) }
//
//func resolveIdent(s Frames, id *Ident) {
//
//	// defScope: where ident was defined
//	// usingScope: where ident is being used:
//	// 	x ->           // defScope of x
//	// 		y ->
//	//  		x + y  // usingScope of x
//
//	name := id.Name
//	v, defScope := s.Find(name)
//	if v == nil {
//		panic(se.Errorf("undefined: %v", name))
//	}
//	//usingScope := s.Last()
//
//	switch {
//	case defScope == len(s)-1: // local variable
//		id.Var = v
//	case defScope == 0: // global variable
//		// TODO
//	default: // captured variable
//		// loop over frames, capture from defscope+1 to last, capture all the way
//		for i := defScope + 1; i < len(s); i++ {
//			v := s[i-1].Find(name)
//			s[i].(*Lambda).DoCapture(name, v.(*LocalVar)) // only locals can be captured
//		}
//		v := s[len(s)-1].Find(name)
//		id.Var = v
//	}
//}
//
//func resolveLambda(s Frames, n *Lambda) {
//	// first define the arguments
//	for i, a := range n.Args {
//		a.Var = &LocalVar{i}
//	}
//
//	// then resolve the body
//	s.Push(n)
//	resolve(s, n.Body)
//	s.Pop()
//}
//
//func (n *Lambda) Find(name string) Var {
//	for _, a := range n.Args {
//		assert(a.Var != nil)
//		if name == a.Name {
//			return a.Var
//		}
//	}
//	for _, a := range n.Cap {
//		if name == a.Name {
//			return a.Local // ?
//		}
//	}
//	return nil
//}
//
//func (n *Lambda) DoCapture(name string, v *LocalVar) (local *CaptVar) {
//	//log.Printf("docapture %q %#v", name, v)
//	if v := n.Find(name); v != nil {
//		return v.(*CaptVar) // already captured
//	}
//	c := &CaptVar{
//		Name:   name,
//		ParVar: v,
//		Local:  &LocalVar{Index: n.NumLocals()},
//	}
//	n.Cap = append(n.Cap, c)
//	return c
//}
//
//func (n *Lambda) NumLocals() int {
//	return len(n.Args) + len(n.Cap)
//}
//
//type Frame interface {
//	Find(name string) Var
//}
//
//type Frames []Frame
//
//func (s *Frames) Push(v Frame) {
//	*s = append(*s, v)
//}
//
//func (s *Frames) Pop() Frame {
//	v := s.Last()
//	*s = (*s)[:len(*s)-1]
//	return v
//}
//
//func (s *Frames) Last() Frame {
//	return (*s)[len(*s)-1]
//}
//
//func (f *Frames) Find(name string) (Var, int) {
//	//defer func() { log.Printf("find %q: %#v", name, v_) }()
//	s := *f
//	for i := len(s) - 1; i >= 0; i-- {
//		if v := s[i].Find(name); v != nil {
//			return v, i
//		}
//	}
//	return nil, 0
//}
