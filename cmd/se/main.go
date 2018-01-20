package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/barnex/se-lang/ast"
)

func main() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		src, err := in.ReadBytes('\n')
		if err != nil {
			return // EOF
		}

		expr, err := ast.Parse(bytes.NewReader(src))
		if err != nil {
			fmt.Println(err)
			continue
		}
		ast.Resolve(expr)

		fmt.Println(ast.ToString(expr))

		//prog, err := eva.CompileAST(expr)
		//if err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		//fmt.Println(eva.Eval(prog))
	}
}
