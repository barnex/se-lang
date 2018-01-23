package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/barnex/se-lang/eva"
)

func main() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		src, err := in.ReadBytes('\n')
		if err != nil {
			return // EOF
		}

		//expr, err := ast.ParseProgram(bytes.NewReader(src))
		//if err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		//ast.Resolve(expr)
		//fmt.Println(ast.ToString(expr))

		prog, err := eva.Compile(bytes.NewReader(src))
		if err != nil {
			fmt.Println(err)
			continue
		}
		v, err := eva.Eval(prog)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%#v\n", v)
		}
	}
}
