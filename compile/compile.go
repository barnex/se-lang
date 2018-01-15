package se

import "io"

func Compile(src io.Reader) (_ Prog, e error) {

	defer func() {
		switch p := recover().(type) {
		default:
			panic(p)
		case nil:
		case *SyntaxError:
			e = p
		}
	}()

	ast := NewParser(src).PExpr()
	Resolve(ast)
	return compileExpr(ast), nil
}

type Prog interface {
	Eval() Value
}

func compileExpr(n Node) Prog {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *Num:
		return &Const{n.Value}
	case *Ident:
		return compileVar(n.Var)
	case *Call:
		return compileCall(n)
	case *Lambda:
		return compileLambda(n)
	}
}
