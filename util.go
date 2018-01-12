package se

import "fmt"

func unhandled(n Node) string {
	return fmt.Sprintf("bug: unhandled AST node: %T", n)
}

func assert(test bool) {
	if !test {
		panic("assertion failed")
	}
}
