package eva

import (
	"fmt"
	"io"

	se "github.com/barnex/se-lang"
	"github.com/barnex/se-lang/ast"
)

func Compile(src io.Reader) (Prog, error) {
	n, err := ast.Parse(src)
	if err != nil {
		return nil, err
	}
	return CompileAST(n)
}

func CompileAST(root ast.Node) (_ Prog, err error) {
	defer func() {
		switch e := recover().(type) {
		case nil: //OK
		default:
			panic(e)
		case se.Error:
			err = e
		}
	}()

	ast.Resolve(root)
	return compileExpr(root), nil
}

func Eval(p Prog) (Value, error) {
	var s Stack
	p.Eval(&s)
	if s.Len() != 1 {
		return nil, fmt.Errorf("got %v values: %v", s.Len, s)
	}
	return s.Pop(), nil
}

func compileExpr(n ast.Node) Prog {
	switch n := n.(type) {
	default:
		panic(unhandled(n))
	case *ast.Call:
		return compileCall(n)
	case *ast.Ident:
		return compileIdent(n)
	case *ast.Lambda:
		return compileLambda(n)
	case *ast.Num:
		return &Const{n.Value}
	}
}

func unhandled(x interface{}) string {
	return fmt.Sprintf("BUG: unhandled case: %T", x)
}
