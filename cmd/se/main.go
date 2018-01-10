package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/barnex/se-lang"
	"github.com/kylelemons/godebug/pretty"
)

func main() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		src, err := in.ReadBytes('\n')
		if err != nil {
			return // EOF
		}
		//expr, err := se.Parse(bytes.NewReader(src))
		//if err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		prog, err := se.Compile(bytes.NewReader(src))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Print(se.ToString(prog), ": ")
		res := prog.Eval()
		fmt.Printf("%T: ", res)
		//if res, ok := res.(interface {
		//	WriteTo(io.Writer)
		//}); ok {
		//	res.WriteTo(os.Stdout)
		//} else {
		//	fmt.Println(res)
		//}
		pretty.CompareConfig.Print(res)
	}
}
