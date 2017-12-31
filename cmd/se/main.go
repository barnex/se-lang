package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/barnex/se-lang"
)

func main() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		src, err := in.ReadBytes('\n')
		if err != nil {
			return // EOF
		}
		expr, err := e.Parse(bytes.NewReader(src))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(e.ToString(expr))
		//eval(expr)
	}
}

//func eval(expr e.Node) {
//	defer func() {
//		if err := recover(); err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	fmt.Println(expr.Eval())
//}
