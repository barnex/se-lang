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
		tree, err := e.Parse(bytes.NewReader(src))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(e.ExprString(tree))
		}
	}
}
