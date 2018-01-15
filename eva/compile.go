package eva

import (
	"fmt"
	"io"

	"github.com/barnex/se-lang/ast"
)

func Compile(src io.Reader) (_ Prog, e error) {
	n, err := ast.Parse(src)
	if err != nil {
		return nil, err
	}

	//defer func() {
	//	switch p := recover().(type) {
	//	default:
	//		panic(p)
	//	case nil:
	//	case se.Error:
	//		e = p
	//	}
	//}()

	return compileExpr(n), nil
}

func compileExpr(n ast.Node) Prog {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *ast.Num:
		return &Const{n.Value}
		//case *ast.Ident:
		//	return compileVar(n.Var)
		//case *ast.Call:
		//	return compileCall(n)
		//case *ast.Lambda:
		//	return compileLambda(n)
	}
}

func unhandled(x interface{}) string {
	return fmt.Sprintf("BUG: unhandled case: %T", x)
}
