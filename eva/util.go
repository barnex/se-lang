package eva

import (
	"fmt"
	"io"

	"github.com/barnex/se-lang/ast"
)

func Eval(p Prog) (Value, error) {
	var m Machine
	p.Exec(&m)
	if len(m.s) != 0 {
		return nil, fmt.Errorf("left dirty stack: %v", m.s)
	}
	return m.RA(), nil
}

func Compile(src io.Reader) (Prog, error) {
	n, err := ast.Parse(src)
	if err != nil {
		return nil, err
	}
	return CompileAST(n)
}

func CompileAST(root ast.Node) (_ Prog, err error) {
	//defer func() {
	//	switch e := recover().(type) {
	//	case nil: //OK
	//	default:
	//		panic(e)
	//	case se.Error:
	//		err = e
	//	}
	//}()

	Resolve(root)
	return compileExpr(root), nil
}

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}
