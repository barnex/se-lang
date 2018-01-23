package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/barnex/se-lang/eva"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if flag.NArg() == 0 {
		repl()
		return
	}

	if flag.NArg() > 1 {
		log.Fatal("too many input files")
		os.Exit(1)
	}

	evalFile(flag.Arg(0))
}

func evalFile(name string) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	prog, err := eva.Compile(bufio.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}
	v, err := eva.Eval(prog)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", v)
}

func repl() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		src, err := in.ReadBytes('\n')
		if err != nil {
			return // EOF
		}

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
